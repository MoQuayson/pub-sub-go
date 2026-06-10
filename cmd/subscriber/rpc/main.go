package main

import (
	"context"
	"github.com/MoQuayson/pub-sub-go/pkg/subscriber"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/constants"
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

	sub, err := subscriber.NewSubscriber(&models.SubscriberConfig{
		Host:               envy.Get("HOST", constants.DefaultHost),
		Port:               envy.Get("PORT", constants.DefaultPort),
		MessagePublishTime: models.LatestPublishTime,
		Partition:          models.DefaultPartition,
		Transport:          models.DefaultTransport,
	})
	if err != nil {
		slog.Error("failed to create subscriber", "error", err)
		os.Exit(1)
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

		slog.Info("received messages", "count", len(messages))
		for _, message := range messages {
			slog.Info("received message", "id", message.Id, "data", message.Data)
		}
	}
}
