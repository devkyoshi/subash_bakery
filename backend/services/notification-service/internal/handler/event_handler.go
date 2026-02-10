package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/yourusername/erp-system/services/notification-service/internal/service"
	"github.com/yourusername/erp-system/shared/rabbitmq"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventHandler struct {
	rabbitClient *rabbitmq.RabbitMQClient
	notifService *service.NotificationService
}

func NewEventHandler(rabbitClient *rabbitmq.RabbitMQClient, notifService *service.NotificationService) *EventHandler {
	return &EventHandler{
		rabbitClient: rabbitClient,
		notifService: notifService,
	}
}

func (h *EventHandler) Start(ctx context.Context) error {
	// Declare queue
	q, err := h.rabbitClient.DeclareQueue("notification_queue")
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Declare exchange
	err = h.rabbitClient.Channel.ExchangeDeclare(
		"notification_events", // name
		"topic",               // type
		true,                  // durable
		false,                 // auto-deleted
		false,                 // internal
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Bind queue to exchange
	err = h.rabbitClient.Channel.QueueBind(
		q.Name,                // queue name
		"inventory.*",         // routing key
		"notification_events", // exchange
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	// Consume
	return h.rabbitClient.Consume(q.Name, h.handleMessage)
}

func (h *EventHandler) handleMessage(body []byte) error {
	var event map[string]interface{}
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	eventType, ok := event["type"].(string)
	if !ok {
		return fmt.Errorf("missing event type")
	}

	orgIDStr, ok := event["organization_id"].(string)
	if !ok {
		return fmt.Errorf("missing organization_id")
	}
	orgID, err := primitive.ObjectIDFromHex(orgIDStr)
	if err != nil {
		return fmt.Errorf("invalid organization_id: %w", err)
	}

	var title, messageBody string

	switch eventType {
	case "low_stock":
		productID := event["product_id"].(string)
		currentStock := event["current_stock"].(float64)
		title = "Low Stock Alert"
		messageBody = fmt.Sprintf("Product %s is running low. Current stock: %.2f", productID, currentStock)
	case "expiry_warning":
		batchNumber := event["batch_number"].(string)
		daysRemaining := event["days_remaining"].(float64)
		title = "Expiry Warning"
		messageBody = fmt.Sprintf("Batch %s expires in %.0f days", batchNumber, daysRemaining)
	default:
		log.Printf("Unknown event type: %s", eventType)
		return nil
	}

	// Convert event map to map[string]string for FCM data payload
	data := make(map[string]string)
	for k, v := range event {
		data[k] = fmt.Sprintf("%v", v)
	}

	// Send notification
	// We use a background context because message handling might be short-lived but sending is network I/O
	// In production, use a timeout context
	ctx := context.Background()
	return h.notifService.SendPushNotification(ctx, orgID, title, messageBody, data)
}
