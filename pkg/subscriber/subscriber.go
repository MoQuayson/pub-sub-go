package subscriber

import (
	"github.com/MoQuayson/pub-sub-go/internal/rpc/subscriber"
	"github.com/MoQuayson/pub-sub-go/pkg/utils"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"log"
)

type Subscriber interface {
	Subscribe(topic string) (models.MessageList, error)
}

//func NewSubscriber(cfg *models.SubscriberConfig) Subscriber {
//	return services.NewSubscriberService(cfg)
//}

func NewSubscriber(cfg *models.SubscriberConfig) Subscriber {

	if cfg == nil {
		return nil
	}

	switch cfg.Transport {
	default:
		return newRpcSubscriber(cfg)
	}
}

func newRpcSubscriber(cfg *models.SubscriberConfig) Subscriber {
	client, err := models.ConnectToRpcServer(cfg.Host, cfg.Port)

	if err != nil {
		log.Fatalf("failed to create subscriber: %v\n", err)
		return nil
	}

	subId := utils.NewSubscriberId()
	if cfg != nil && cfg.SubscriberId != nil {
		subId = *cfg.SubscriberId
	}

	return &subscriber.RpcSubscriber{
		Client: client,
		Id:     subId,
		Config: cfg,
	}
}
