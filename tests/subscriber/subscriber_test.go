package publisher

import (
	"github.com/MoQuayson/pub-sub-go/pkg/subscriber"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"log"
	"testing"
)

func TestSubscriberService(t *testing.T) {
	sub := subscriber.NewSubscriber(&models.SubscriberConfig{
		Host:               "0.0.0.0",
		Port:               "7000",
		MessagePublishTime: models.LatestPublishTime,
	})

	for {
		messages, err := sub.Subscribe("test")
		if err != nil {
			t.Errorf("failed to publish message: %v\n", err)
		}

		for _, message := range messages {
			log.Printf("Message (%s) with data:%v\n", message.Id, message.Data)
		}
	}
}
