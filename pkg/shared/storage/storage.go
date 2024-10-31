package storage

import (
	"github.com/MoQuayson/go-event-bridge/pkg/shared/models"
	"time"
)

// Storage defines methods to be implemented by each storage backend.
type Storage interface {
	StoreMessage(msg *models.Message) error
	GetMessages(topic string, partition models.Partition) ([]*models.Message, error)
	EvictMessages(topic string, partition models.Partition, ttl time.Duration) error
}
