package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQClient struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	URL     string
	mu      sync.Mutex
}

func NewRabbitMQClient(url string) (*RabbitMQClient, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	return &RabbitMQClient{
		Conn:    conn,
		Channel: ch,
		URL:     url,
	}, nil
}

func (c *RabbitMQClient) Close() {
	if c.Channel != nil {
		c.Channel.Close()
	}
	if c.Conn != nil {
		c.Conn.Close()
	}
}

func (c *RabbitMQClient) Publish(ctx context.Context, exchange, routingKey string, body interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	bytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal body: %w", err)
	}

	return c.Channel.PublishWithContext(ctx,
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bytes,
			Timestamp:   time.Now(),
		})
}

func (c *RabbitMQClient) Consume(queueName string, handler func([]byte) error) error {
	msgs, err := c.Channel.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	go func() {
		for d := range msgs {
			if err := handler(d.Body); err != nil {
				log.Printf("Error handling message: %v", err)
			}
		}
	}()

	return nil
}

func (c *RabbitMQClient) DeclareQueue(name string) (amqp.Queue, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Channel.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
}

// RPCClient sends a request and waits for a response
func (c *RabbitMQClient) RPCRequest(ctx context.Context, queueName string, request interface{}) (interface{}, error) {
	ch, err := c.Conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name (empty means random)
		false, // durable
		true,  // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare reply queue: %w", err)
	}
	log.Printf("RPC Client declared reply queue: %s", q.Name)

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register consumer: %w", err)
	}

	corrId := randomString(32)

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	err = ch.PublishWithContext(ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          body,
		})
	if err != nil {
		return nil, fmt.Errorf("failed to publish request: %w", err)
	}

	select {
	case d := <-msgs:
		if d.CorrelationId == corrId {
			var response interface{}
			if err := json.Unmarshal(d.Body, &response); err != nil {
				return nil, fmt.Errorf("failed to unmarshal response: %w", err)
			}
			return response, nil
		}
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(5 * time.Second):
		return nil, fmt.Errorf("timeout waiting for response")
	}
	return nil, fmt.Errorf("unknown error")
}

// RPCServer listens for requests and sends responses
func (c *RabbitMQClient) RPCServe(queueName string, handler func([]byte) (interface{}, error)) error {
	// Use a dedicated channel for the consumer to avoid concurrency issues
	ch, err := c.Conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel for RPC server: %w", err)
	}

	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	go func() {
		defer ch.Close()
		for d := range msgs {
			log.Printf("RPC Server received message on queue %s with CorrelationId %s, ReplyTo: %s", queueName, d.CorrelationId, d.ReplyTo)
			response, err := handler(d.Body)
			if err != nil {
				log.Printf("Error handling RPC request: %v", err)
				// Look into sending error response
			}

			responseBody, _ := json.Marshal(response)

			log.Printf("RPC Server sending response to queue %s with CorrelationId %s", d.ReplyTo, d.CorrelationId)
			err = ch.PublishWithContext(context.Background(),
				"",        // exchange
				d.ReplyTo, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType:   "application/json",
					CorrelationId: d.CorrelationId,
					Body:          responseBody,
				})
			if err != nil {
				log.Printf("Failed to publish response: %v", err)
			}

			d.Ack(false)
		}
	}()

	return nil
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(65 + i%26)
	}
	return string(bytes)
}
