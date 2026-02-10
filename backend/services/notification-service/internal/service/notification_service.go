package service

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/yourusername/erp-system/services/notification-service/config"
	"github.com/yourusername/erp-system/services/notification-service/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/api/option"
)

type NotificationService struct {
	deviceRepo *repository.DeviceRepository
	fcmClient  *messaging.Client
}

func NewNotificationService(cfg *config.Config, deviceRepo *repository.DeviceRepository) (*NotificationService, error) {
	// Initialize Firebase
	opt := option.WithCredentialsFile(cfg.FirebaseCredentialsPath)
	conf := &firebase.Config{ProjectID: cfg.FirebaseProjectID}
	app, err := firebase.NewApp(context.Background(), conf, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error getting Messaging client: %v", err)
	}

	return &NotificationService{
		deviceRepo: deviceRepo,
		fcmClient:  client,
	}, nil
}

// SendPushNotification sends a notification to all devices in an organization
func (s *NotificationService) SendPushNotification(ctx context.Context, orgID primitive.ObjectID, title, body string, data map[string]string) error {
	tokens, err := s.deviceRepo.FindByOrganizationID(ctx, orgID)
	if err != nil {
		return fmt.Errorf("failed to fetch device tokens: %w", err)
	}

	if len(tokens) == 0 {
		log.Printf("No devices found for organization %s", orgID.Hex())
		return nil
	}

	// FCM Multicast allows up to 500 tokens per message
	batchSize := 500
	for i := 0; i < len(tokens); i += batchSize {
		end := i + batchSize
		if end > len(tokens) {
			end = len(tokens)
		}
		batchTokens := tokens[i:end]

		message := &messaging.MulticastMessage{
			Tokens: batchTokens,
			Notification: &messaging.Notification{
				Title: title,
				Body:  body,
			},
			Data: data,
		}

		br, err := s.fcmClient.SendMulticast(ctx, message)
		if err != nil {
			log.Printf("Error sending batch: %v", err)
			continue
		}

		if br.FailureCount > 0 {
			var failedTokens []string
			for idx, resp := range br.Responses {
				if !resp.Success {
					if resp.Error != nil && messaging.IsRegistrationTokenNotRegistered(resp.Error) {
						failedTokens = append(failedTokens, batchTokens[idx])
					}
				}
			}

			if len(failedTokens) > 0 {
				log.Printf("Cleaning up %d invalid tokens", len(failedTokens))
				for _, t := range failedTokens {
					_ = s.deviceRepo.Delete(ctx, t)
				}
			}
		}
		log.Printf("Sent notification batch: %d success, %d failure", br.SuccessCount, br.FailureCount)
	}

	return nil
}
