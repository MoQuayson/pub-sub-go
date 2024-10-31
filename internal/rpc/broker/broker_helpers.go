package broker

import (
	"github.com/MoQuayson/go-event-bridge/pkg/shared/models"
	storage2 "github.com/MoQuayson/go-event-bridge/pkg/shared/storage"
	"github.com/MoQuayson/go-event-bridge/pkg/shared/utils/enums"
)

func getStorage(config *models.BrokerConfig) storage2.Storage {
	storageType := config.Storage
	switch storageType {
	case enums.StorageType_InMemory:
		return storage2.NewInMemoryStorage()
	case enums.StorageType_Redis:
		return nil
	case enums.StorageType_Disk:
		return nil
	default:
		return storage2.NewInMemoryStorage()
	}
}
