package publisher

import (
	"github.com/MoQuayson/pub-sub-go/internal/rpc/publisher"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"log"
)

type Publisher interface {
	PublishMessage(topic string, partition models.Partition, data string) error
}

func NewPublisher(cfg *models.PublisherConfig) Publisher {

	if cfg == nil {
		return nil
	}

	switch cfg.Transport {
	default:
		return newRpcPublisher(cfg)
	}
}

func newRpcPublisher(cfg *models.PublisherConfig) Publisher {
	pub := &publisher.RpcPublisher{}
	client, err := models.ConnectToRpcServer(cfg.Host, cfg.Port)

	if err != nil {
		log.Fatalf("failed to create a publisher: %v\n", err)
		return nil
	}

	pub.Client = client
	return pub
}
