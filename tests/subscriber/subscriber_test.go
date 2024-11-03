package publisher

import (
	"github.com/MoQuayson/pub-sub-go/pkg/shared/models"
	"github.com/MoQuayson/pub-sub-go/pkg/subscriber"
	"github.com/MoQuayson/pub-sub-go/tests/shared"
	"log"
	"testing"
	"time"
)

func TestSubscriberService(t *testing.T) {
	sub := subscriber.NewSubscriber(shared.ConnectionConfig)

	for {
		messages, err := sub.GetMessages("test", models.DefaultPartition, time.Now().Add(-1))
		if err != nil {
			t.Errorf("failed to publish message: %v\n", err)
		}

		for _, message := range messages {
			log.Printf("Message (%s) with data:%v\n", message.Id, message.Data)
		}
	}
}
