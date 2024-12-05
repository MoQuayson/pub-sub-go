package main

import (
	"github.com/MoQuayson/pub-sub-go/pkg/broker"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"github.com/gobuffalo/envy"
	"log"
)

func main() {
	grpcBroker := broker.NewBroker(&models.BrokerConfig{
		Host:      envy.Get("HOST", "0.0.0.0"),
		Port:      envy.Get("PORT", "50051"),
		Transport: models.GrpcTransport,
		Storage:   models.InMemoryStorageType,
	})

	log.Fatalln(grpcBroker.Start())
}
