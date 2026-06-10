package main

import (
	"context"
	"github.com/MoQuayson/pub-sub-go/pkg/broker"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"github.com/gobuffalo/envy"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	storageType := envy.Get("STORAGE", "InMemory")
	var (
		b   broker.Broker
		err error
	)
	if storageType == "Disk" {
		b, err = getBrokerWithDiskStorage()
	} else {
		b, err = getBrokerWithInMemoryStorage()
	}
	if err != nil {
		slog.Error("failed to create broker", "error", err)
		os.Exit(1)
	}

	if err := b.Start(ctx); err != nil {
		slog.Error("broker stopped with error", "error", err)
		os.Exit(1)
	}
	slog.Info("broker stopped cleanly")
}

func getBrokerWithInMemoryStorage() (broker.Broker, error) {
	return broker.NewBroker(&models.BrokerConfig{
		Host:      envy.Get("HOST", "127.0.0.1"),
		Port:      envy.Get("PORT", "7000"),
		Transport: models.DefaultTransport,
		Storage:   models.InMemoryStorageType,
	})
}

func getBrokerWithDiskStorage() (broker.Broker, error) {
	loc := envy.Get("DISK_STORAGE", "./bin/logs")
	return broker.NewBroker(&models.BrokerConfig{
		Host:            envy.Get("HOST", "127.0.0.1"),
		Port:            envy.Get("PORT", "7000"),
		Transport:       models.DefaultTransport,
		Storage:         models.DiskStorageType,
		StorageLocation: &loc,
	})
}
