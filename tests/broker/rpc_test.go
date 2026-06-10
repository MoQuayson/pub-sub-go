package broker

import (
	"context"
	"github.com/MoQuayson/pub-sub-go/pkg/broker"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"testing"
	"time"
)

const logDir = "../../bin/logs"

func TestRpcBrokerServerWithInMemoryStorage(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)

	brokerSrv, err := broker.NewBroker(&models.BrokerConfig{
		Host:      "127.0.0.1",
		Port:      "17001",
		Storage:   models.InMemoryStorageType,
		Transport: models.DefaultTransport,
	})
	if err != nil {
		t.Fatal(err)
	}

	startErr := make(chan error, 1)
	go func() { startErr <- brokerSrv.Start(ctx) }()

	time.Sleep(50 * time.Millisecond)

	msg := &models.Message{
		Id:        "MSG001",
		Topic:     "test",
		Partition: models.DefaultPartition,
		Data:      "hello",
		Timestamp: time.Now(),
	}
	if err := brokerSrv.PublishMessage(msg); err != nil {
		t.Fatalf("publish failed: %v", err)
	}

	cancel()
	select {
	case err := <-startErr:
		if err != nil {
			t.Errorf("Start returned error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Error("server did not stop within 2s")
	}
}

func TestRpcBrokerServerWithDiskStorage(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)

	loc := t.TempDir()
	brokerSrv, err := broker.NewBroker(&models.BrokerConfig{
		Host:            "127.0.0.1",
		Port:            "17002",
		Storage:         models.DiskStorageType,
		StorageLocation: &loc,
		Transport:       models.DefaultTransport,
	})
	if err != nil {
		t.Fatal(err)
	}

	startErr := make(chan error, 1)
	go func() { startErr <- brokerSrv.Start(ctx) }()

	time.Sleep(50 * time.Millisecond)

	msg := &models.Message{
		Id:        "MSG002",
		Topic:     "test",
		Partition: models.DefaultPartition,
		Data:      "hello disk",
		Timestamp: time.Now(),
	}
	if err := brokerSrv.PublishMessage(msg); err != nil {
		t.Fatalf("publish failed: %v", err)
	}

	cancel()
	select {
	case err := <-startErr:
		if err != nil {
			t.Errorf("Start returned error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Error("server did not stop within 2s")
	}
}
