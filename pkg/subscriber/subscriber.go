package subscriber

import (
	"github.com/MoQuayson/pub-sub-go/internal/subscriber/services"
	"github.com/MoQuayson/pub-sub-go/pkg/shared/models"
)

type Subscriber interface {
	Subscribe(topic string) (models.MessageList, error)
}

func NewSubscriber(cfg *models.SubscriberConfig) Subscriber {
	return services.NewSubscriberService(cfg)
}
