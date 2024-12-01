package main

import (
	"fmt"
	"github.com/MoQuayson/pub-sub-go/pkg/publisher"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"github.com/gobuffalo/envy"
	"log"
)

func main() {
	pub := publisher.NewPublisher(&models.PublisherConfig{
		Host:      envy.Get("HOST", ""),
		Port:      envy.Get("PORT", ""),
		Transport: models.DefaultTransport,
	})

	if err := pub.PublishMessage("test", models.DefaultPartition, "Testing Data"); err != nil {
		log.Fatalln("failed to publish message: ", err)
	}

	//multiple publishing
	for i := 1; i <= 9_000; i++ {
		if err := pub.PublishMessage("test", models.DefaultPartition, fmt.Sprintf("Data %d", i)); err != nil {
			log.Fatalln("failed to publish message: ", err)
		}
	}

	fmt.Println("published successfully")
}
