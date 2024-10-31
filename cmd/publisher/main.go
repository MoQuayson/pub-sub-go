package main

import (
	"fmt"
	"github.com/MoQuayson/go-event-bridge/pkg/shared/models"
	"github.com/MoQuayson/go-event-bridge/publisher"
	"github.com/gobuffalo/envy"
	"log"
)

func main() {
	pub := publisher.NewPublisher(&models.RpcConnConfig{
		Host: envy.Get("HOST", ""),
		Port: envy.Get("PORT", ""),
	})

	//time.Sleep(5 * time.Second)
	//single publish
	if err := pub.PublishMessage("test", models.DefaultPartition, "Testing Data"); err != nil {
		log.Fatalln("failed to publish message: ", err)
	}

	//multiple publishing
	for i := 1; i <= 9_000; i++ {
		if err := pub.PublishMessage("test", models.DefaultPartition, fmt.Sprintf("EData %d", i)); err != nil {
			log.Fatalln("failed to publish message: ", err)
		}
	}
}
