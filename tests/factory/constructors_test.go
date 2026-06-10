package factory

import (
	"testing"

	"github.com/MoQuayson/pub-sub-go/pkg/broker"
	"github.com/MoQuayson/pub-sub-go/pkg/publisher"
	"github.com/MoQuayson/pub-sub-go/pkg/storage"
	"github.com/MoQuayson/pub-sub-go/pkg/subscriber"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
)

func TestConstructorsRejectNilConfig(t *testing.T) {
	if _, err := broker.NewBroker(nil); err == nil {
		t.Fatal("expected broker nil config error")
	}
	if _, err := publisher.NewPublisher(nil); err == nil {
		t.Fatal("expected publisher nil config error")
	}
	if _, err := subscriber.NewSubscriber(nil); err == nil {
		t.Fatal("expected subscriber nil config error")
	}
}

func TestConstructorsRejectUnsupportedTransport(t *testing.T) {
	if _, err := broker.NewBroker(&models.BrokerConfig{Transport: models.WebSocketTransport}); err == nil {
		t.Fatal("expected broker unsupported transport error")
	}
	if _, err := publisher.NewPublisher(&models.PublisherConfig{Transport: models.WebSocketTransport}); err == nil {
		t.Fatal("expected publisher unsupported transport error")
	}
	if _, err := subscriber.NewSubscriber(&models.SubscriberConfig{Transport: models.WebSocketTransport}); err == nil {
		t.Fatal("expected subscriber unsupported transport error")
	}
}

func TestStorageRejectsUnsupportedBackends(t *testing.T) {
	if _, err := storage.GetStorage(&models.BrokerConfig{Storage: models.RedisStorageType}); err == nil {
		t.Fatal("expected redis storage error")
	}
	if _, err := storage.GetStorage(&models.BrokerConfig{Storage: models.DiskStorageType}); err == nil {
		t.Fatal("expected missing disk location error")
	}
}

func TestStorageDefaultsToInMemory(t *testing.T) {
	if _, err := storage.GetStorage(&models.BrokerConfig{}); err != nil {
		t.Fatalf("expected empty storage type to default to in-memory: %v", err)
	}
}
