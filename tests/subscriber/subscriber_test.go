package publisher

import (
	"context"
	"github.com/MoQuayson/pub-sub-go/pkg/broker"
	"github.com/MoQuayson/pub-sub-go/pkg/publisher"
	"github.com/MoQuayson/pub-sub-go/pkg/subscriber"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"io"
	"testing"
	"time"
)

func startBrokerAndPublish(t *testing.T, port string, transport models.Transport, data string) context.CancelFunc {
	t.Helper()
	b, err := broker.NewBroker(&models.BrokerConfig{
		Host:      "127.0.0.1",
		Port:      port,
		Storage:   models.InMemoryStorageType,
		Transport: transport,
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	go b.Start(ctx)
	time.Sleep(50 * time.Millisecond)

	pub, err := publisher.NewPublisher(&models.PublisherConfig{
		Host:      "127.0.0.1",
		Port:      port,
		Transport: transport,
	})
	if err != nil {
		cancel()
		t.Fatal(err)
	}
	if err := pub.PublishMessage(context.Background(), "test", models.DefaultPartition, data); err != nil {
		cancel()
		t.Fatalf("setup publish failed: %v", err)
	}
	return cancel
}

func TestRpcSubscriberService(t *testing.T) {
	t.Parallel()
	cancel := startBrokerAndPublish(t, "17020", models.DefaultTransport, "hello rpc")
	t.Cleanup(cancel)

	ctx := context.Background()
	sub, err := subscriber.NewSubscriber(&models.SubscriberConfig{
		Host:               "127.0.0.1",
		Port:               "17020",
		MessagePublishTime: models.EarliestPublishTime,
		Partition:          models.DefaultPartition,
		Transport:          models.DefaultTransport,
	})
	if err != nil {
		t.Fatal(err)
	}

	var got []*models.Message
	for i := 0; i < 5; i++ {
		got, err = sub.Subscribe(ctx, "test")
		if err != nil {
			t.Fatalf("Subscribe: %v", err)
		}
		if len(got) > 0 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	if len(got) == 0 {
		t.Error("expected at least one message, got none")
	}
}

func TestGrpcSubscriberService(t *testing.T) {
	t.Parallel()
	cancel := startBrokerAndPublish(t, "17021", models.GrpcTransport, "hello grpc")
	t.Cleanup(cancel)

	ctx := context.Background()
	sub, err := subscriber.NewSubscriber(&models.SubscriberConfig{
		Host:               "127.0.0.1",
		Port:               "17021",
		MessagePublishTime: models.EarliestPublishTime,
		Partition:          models.DefaultPartition,
		Transport:          models.GrpcTransport,
	})
	if err != nil {
		t.Fatal(err)
	}
	if c, ok := sub.(io.Closer); ok {
		t.Cleanup(func() { c.Close() })
	}

	var got []*models.Message
	for i := 0; i < 5; i++ {
		got, err = sub.Subscribe(ctx, "test")
		if err != nil {
			t.Fatalf("Subscribe: %v", err)
		}
		if len(got) > 0 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	if len(got) == 0 {
		t.Error("expected at least one message, got none")
	}
}
