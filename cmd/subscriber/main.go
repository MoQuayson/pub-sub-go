package main

import (
	"github.com/MoQuayson/go-event-bridge/pkg/shared/models"
	"github.com/MoQuayson/go-event-bridge/subscriber"
	"github.com/gobuffalo/envy"
	"log"
	"time"
)

func main() {
	sub := subscriber.NewSubscriber(&models.RpcConnConfig{
		Host: envy.Get("HOST", ""),
		Port: envy.Get("PORT", ""),
	})

	for {
		messages, err := sub.GetMessages("test", models.DefaultPartition, time.Now().Add(-1))
		if err != nil {
			log.Fatalf("failed to publish message: %v\n", err)
		}

		for _, message := range messages {
			log.Printf("Message (%s) with data:%v\n", message.Id, message.Data)
		}
	}

}
