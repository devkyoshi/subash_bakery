package worker

import (
	"context"
	"log"
	"time"

	"github.com/yourusername/erp-system/services/inventory-service/internal/repository"
	"github.com/yourusername/erp-system/shared/rabbitmq"
)

type ExpiryChecker struct {
	batchRepo    *repository.BatchRepository
	rabbitClient *rabbitmq.RabbitMQClient
	daysWarning  int
}

func NewExpiryChecker(batchRepo *repository.BatchRepository, rabbitClient *rabbitmq.RabbitMQClient, daysWarning int) *ExpiryChecker {
	return &ExpiryChecker{
		batchRepo:    batchRepo,
		rabbitClient: rabbitClient,
		daysWarning:  daysWarning,
	}
}

func (w *ExpiryChecker) Start(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	// Run once immediately on startup
	w.CheckExpiry(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.CheckExpiry(ctx)
		}
	}
}

func (w *ExpiryChecker) CheckExpiry(ctx context.Context) {
	log.Println("Running expiry check...")
	batches, err := w.batchRepo.FindExpiringBatches(ctx, w.daysWarning)
	if err != nil {
		log.Printf("Failed to find expiring batches: %v", err)
		return
	}

	for _, batch := range batches {
		event := map[string]interface{}{
			"type":            "expiry_warning",
			"batch_id":        batch.ID.Hex(),
			"product_id":      batch.ProductID.Hex(),
			"batch_number":    batch.BatchNumber,
			"expiry_date":     batch.ExpiryDate,
			"organization_id": batch.OrganizationID.Hex(),
			"days_remaining":  int(time.Until(*batch.ExpiryDate).Hours() / 24),
		}

		// Use a dedicated exchange for notifications
		err := w.rabbitClient.Publish(ctx, "notification_events", "inventory.expiry", event)
		if err != nil {
			log.Printf("Failed to publish expiry warning for batch %s: %v", batch.BatchNumber, err)
		}
	}
	log.Printf("Expiry check completed. Found %d expiring batches.", len(batches))
}
