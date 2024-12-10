package publisher

import (
	"github.com/MoQuayson/pub-sub-go/pkg/publisher"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"testing"
)

func TestPublisherService(t *testing.T) {
	pub := publisher.NewPublisher(&models.PublisherConfig{
		Host:      "0.0.0.0",
		Port:      "7000",
		Transport: models.DefaultTransport,
	})

	if err := pub.PublishMessage("test", models.DefaultPartition, "Hello world"); err != nil {
		t.Errorf("failed to publish message: %v\n", err)
	}
	t.Log("publisher service working successfully")
}

func TestGrpcPublisherService(t *testing.T) {
	pub := publisher.NewPublisher(&models.PublisherConfig{
		Host:      "0.0.0.0",
		Port:      "7000",
		Transport: models.GrpcTransport,
	})

	if err := pub.PublishMessage("test", models.DefaultPartition, "Hello world"); err != nil {
		t.Errorf("failed to publish message: %v\n", err)
	}
	t.Log("publisher service working successfully")
}
