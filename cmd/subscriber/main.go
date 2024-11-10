package main

import (
	"github.com/MoQuayson/pub-sub-go/pkg/shared/models"
	"github.com/MoQuayson/pub-sub-go/pkg/subscriber"
	"github.com/gobuffalo/envy"
	"log"
	"time"
)

func main() {
	sub := subscriber.NewSubscriber(&models.BrokerConfig{
		Host: envy.Get("HOST", ""),
		Port: envy.Get("PORT", ""),
	})

	for {
		messages, err := sub.GetMessages("test", models.DefaultPartition, time.Now().Add(-1))
		if err != nil {
			log.Fatalf("failed to subscribe message: %v\n", err)
		}

		for _, message := range messages {
			log.Printf("Message (%s) with data:%v\n", message.Id, message.Data)
		}
	}

}
