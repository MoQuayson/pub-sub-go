package main

import (
	"context"
	"github.com/MoQuayson/pub-sub-go/pkg/subscriber"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/constants"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"github.com/gobuffalo/envy"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	sub, err := subscriber.NewSubscriber(&models.SubscriberConfig{
		Host:               envy.Get("HOST", constants.DefaultHost),
		Port:               envy.Get("GRPC_PORT", "50051"),
		MessagePublishTime: models.EarliestPublishTime,
		Partition:          models.DefaultPartition,
		Transport:          models.GrpcTransport,
	})
	if err != nil {
		slog.Error("failed to create subscriber", "error", err)
		os.Exit(1)
	}
	if c, ok := sub.(io.Closer); ok {
		defer c.Close()
	}

	for {
		select {
		case <-ctx.Done():
			slog.Info("subscriber stopping")
			return
		default:
		}
		messages, err := sub.Subscribe(ctx, "test")
		if err != nil {
			slog.Error("failed to subscribe", "error", err)
			os.Exit(1)
		}
		for _, message := range messages {
			slog.Info("received message", "id", message.Id, "data", message.Data)
		}
	}
}
