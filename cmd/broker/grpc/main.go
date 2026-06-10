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

	grpcBroker, err := broker.NewBroker(&models.BrokerConfig{
		Host:      envy.Get("HOST", "0.0.0.0"),
		Port:      envy.Get("GRPC_PORT", "50051"),
		Transport: models.GrpcTransport,
		Storage:   models.InMemoryStorageType,
	})
	if err != nil {
		slog.Error("failed to create broker", "error", err)
		os.Exit(1)
	}

	if err := grpcBroker.Start(ctx); err != nil {
		slog.Error("broker stopped with error", "error", err)
		os.Exit(1)
	}
	slog.Info("broker stopped cleanly")
}
