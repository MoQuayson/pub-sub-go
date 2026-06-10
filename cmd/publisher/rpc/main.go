package main

import (
	"context"
	"fmt"
	"github.com/MoQuayson/pub-sub-go/pkg/publisher"
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

	pub, err := publisher.NewPublisher(&models.PublisherConfig{
		Host:      envy.Get("HOST", constants.DefaultHost),
		Port:      envy.Get("PORT", constants.DefaultPort),
		Transport: models.DefaultTransport,
	})
	if err != nil {
		slog.Error("failed to create publisher", "error", err)
		os.Exit(1)
	}

	for i := 1; i <= 9_000; i++ {
		if err := pub.PublishMessage(ctx, "test", models.DefaultPartition, fmt.Sprintf("Data %05d", i)); err != nil {
			slog.Error("failed to publish message", "error", err)
			os.Exit(1)
		}
	}

	slog.Info("published successfully")
}
