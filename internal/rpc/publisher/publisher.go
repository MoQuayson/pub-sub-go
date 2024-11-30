package publisher

import (
	"fmt"
	"github.com/MoQuayson/pub-sub-go/pkg/utils"
	constants2 "github.com/MoQuayson/pub-sub-go/pkg/utils/constants"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"log"
	"net/rpc"
	"time"
)

type RpcPublisher struct {
	Client *rpc.Client
}

func NewPublisherService(cfg *models.PublisherConfig) *RpcPublisher {
	pub := &RpcPublisher{}
	Client, err := connectToRpcServer(cfg)

	if err != nil {
		log.Fatalf("failed to create a publisher: %v\n", err)
		return nil
	}

	pub.Client = Client
	return pub
}

// PublishMessage publishes to a topic
func (p *RpcPublisher) PublishMessage(topic string, partition models.Partition, data string) error {
	msg := models.Message{
		Id:        utils.NewMessageId(),
		Topic:     topic,
		Partition: partition,
		Data:      data,
		Timestamp: time.Now(),
	}

	//make an rpc call to a service method
	var reply string
	if err := p.Client.Call(constants2.PublishToBrokerServiceMethod, msg, &reply); err != nil {
		return err
	}

	return nil

}

func connectToRpcServer(c *models.PublisherConfig) (*rpc.Client, error) {
	if c == nil {
		*c = models.PublisherConfig{Host: constants2.DefaultHost, Port: constants2.DefaultPort}
	}
	return rpc.Dial("tcp", fmt.Sprintf("%s:%s", c.Host, c.Port))
}
