package storage

import (
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"time"
)

// Storage defines methods to be implemented by each storage backend.
type Storage interface {
	StoreMessage(msg *models.Message) error
	GetMessages(topic string, partition models.Partition) (models.MessageList, error)
	EvictMessages(topic string, partition models.Partition, ttl time.Duration) error
}

func GetStorage(config *models.BrokerConfig) Storage {
	storageType := config.Storage
	switch storageType {
	case models.InMemoryStorageType:
		return NewInMemoryStorage()
	case models.RedisStorageType:
		return nil
	case models.DiskStorageType:
		return nil
	default:
		return NewInMemoryStorage()
	}
}
