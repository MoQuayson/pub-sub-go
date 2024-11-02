package publisher

import (
	"github.com/MoQuayson/pub-sub-go/internal/publisher/services"
	"github.com/MoQuayson/pub-sub-go/pkg/shared/models"
)

type Publisher interface {
	PublishMessage(topic string, partition models.Partition, data string) error
}

func NewPublisherService(cfg *models.RpcConnConfig) Publisher {
	return services.NewPublisherService(cfg)
}
