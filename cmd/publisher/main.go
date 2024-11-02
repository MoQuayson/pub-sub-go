package main

import (
	"fmt"
	"github.com/MoQuayson/pub-sub-go/pkg/publisher"
	"github.com/MoQuayson/pub-sub-go/pkg/shared/models"
	"github.com/gobuffalo/envy"
	"log"
)

func main() {
	pub := publisher.NewPublisherService(&models.RpcConnConfig{
		Host: envy.Get("HOST", ""),
		Port: envy.Get("PORT", ""),
	})

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
