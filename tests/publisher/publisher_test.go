package publisher

import (
	"github.com/MoQuayson/pub-sub-go/pkg/publisher"
	"github.com/MoQuayson/pub-sub-go/pkg/shared/models"
	"github.com/MoQuayson/pub-sub-go/tests/shared"
	"testing"
)

func TestPublisherService(t *testing.T) {
	pub := publisher.NewPublisher(shared.ConnectionConfig)

	if err := pub.PublishMessage("test", models.DefaultPartition, "Hello world"); err != nil {
		t.Errorf("failed to publish message: %v\n", err)
	}
	t.Log("publisher service working successfully")
}
