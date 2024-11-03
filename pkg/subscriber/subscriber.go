package subscriber

import (
	"github.com/MoQuayson/pub-sub-go/internal/subscriber/services"
	"github.com/MoQuayson/pub-sub-go/pkg/shared/models"
	"time"
)

type Subscriber interface {
	GetMessages(topic string, partition models.Partition, startTime time.Time) (models.MessageList, error)
}

func NewSubscriber(cfg *models.BrokerConfig) Subscriber {
	return services.NewSubscriberService(cfg)
}
