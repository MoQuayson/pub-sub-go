package storage

import (
	"fmt"
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

func GetStorage(config *models.BrokerConfig) (Storage, error) {
	if config == nil {
		return nil, fmt.Errorf("storage: broker config is required")
	}

	storageType := config.Storage
	switch storageType {
	case "", models.InMemoryStorageType:
		return NewInMemoryStorage(), nil
	case models.RedisStorageType:
		return nil, fmt.Errorf("storage: %s storage is not implemented", storageType)
	case models.DiskStorageType:
		if config.StorageLocation == nil || *config.StorageLocation == "" {
			return nil, fmt.Errorf("storage: disk storage location is required")
		}
		w, err := disk.NewLogWriter(*config.StorageLocation)
		if err != nil {
			return nil, err
		}
		return disk.NewDiskStorage(w, disk.NewLogReader(*config.StorageLocation)), nil
	default:
		return nil, fmt.Errorf("storage: unsupported storage type %q", storageType)
	}
}
