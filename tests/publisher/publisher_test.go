package publisher

import (
	"context"
	"github.com/MoQuayson/pub-sub-go/pkg/broker"
	"github.com/MoQuayson/pub-sub-go/pkg/publisher"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"io"
	"testing"
	"time"
)

func startTestBroker(t *testing.T, port string) context.CancelFunc {
	t.Helper()
	b, err := broker.NewBroker(&models.BrokerConfig{
		Host:      "127.0.0.1",
		Port:      port,
		Storage:   models.InMemoryStorageType,
		Transport: models.DefaultTransport,
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	go b.Start(ctx)
	time.Sleep(50 * time.Millisecond)
	return cancel
}

func TestPublisherService(t *testing.T) {
	t.Parallel()
	cancel := startTestBroker(t, "17010")
	t.Cleanup(cancel)

	ctx := context.Background()
	pub, err := publisher.NewPublisher(&models.PublisherConfig{
		Host:      "127.0.0.1",
		Port:      "17010",
		Transport: models.DefaultTransport,
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := pub.PublishMessage(ctx, "test", models.DefaultPartition, "Hello world"); err != nil {
		t.Errorf("PublishMessage: %v", err)
	}
}

func TestGrpcPublisherService(t *testing.T) {
	t.Parallel()

	b, err := broker.NewBroker(&models.BrokerConfig{
		Host:      "127.0.0.1",
		Port:      "17011",
		Storage:   models.InMemoryStorageType,
		Transport: models.GrpcTransport,
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go b.Start(ctx)
	time.Sleep(50 * time.Millisecond)

	pub, err := publisher.NewPublisher(&models.PublisherConfig{
		Host:      "127.0.0.1",
		Port:      "17011",
		Transport: models.GrpcTransport,
	})
	if err != nil {
		t.Fatal(err)
	}
	if c, ok := pub.(io.Closer); ok {
		t.Cleanup(func() { c.Close() })
	}

	if err := pub.PublishMessage(ctx, "test", models.DefaultPartition, "Hello grpc"); err != nil {
		t.Errorf("PublishMessage: %v", err)
	}
}
