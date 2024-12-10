package storage

import (
	"github.com/MoQuayson/pub-sub-go/pkg/storage/disk"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"time"
)

// Storage defines methods to be implemented by each storage backend.
type Storage interface {
	//StoreMessage saves message to a storage
	StoreMessage(msg *models.Message) error
	//GetMessages retrieves messages from a storage
	GetMessages(topic string, partition models.Partition) (models.MessageList, error)

	//EvictMessages removes old messages from storage
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
		return disk.NewDiskStorage(disk.NewLogWriter(*config.StorageLocation), disk.NewLogReader(*config.StorageLocation))
	default:
		return NewInMemoryStorage()
	}
}
