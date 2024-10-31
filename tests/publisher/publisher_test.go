package publisher

import (
	"github.com/MoQuayson/go-event-bridge/pkg/shared/models"
	"github.com/MoQuayson/go-event-bridge/publisher"
	"github.com/MoQuayson/go-event-bridge/tests/shared"
	"testing"
)

func TestPublisherService(t *testing.T) {
	pub := publisher.NewPublisher(shared.ConnectionConfig)

	if err := pub.PublishMessage("test", models.DefaultPartition, "Hello world"); err != nil {
		t.Errorf("failed to publish message: %v\n", err)
	}
	t.Log("publisher service working successfully")
}
