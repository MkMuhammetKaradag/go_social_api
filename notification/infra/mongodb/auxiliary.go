package mongodb

import (
	"context"
	"log"
	"time"
)

func (r *Repository) StartNotificationCleanupTask(interval time.Duration, olderThan time.Duration) {
	go func() {
		for {
			err := r.SoftDeleteOldNotifications(context.Background(), olderThan)
			if err != nil {
				log.Printf("Error during soft delete: %v", err)
			}
			time.Sleep(interval)
		}
	}()
}
