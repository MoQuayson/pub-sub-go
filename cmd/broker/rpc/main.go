package main

import (
	"github.com/MoQuayson/pub-sub-go/pkg/broker"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"github.com/gobuffalo/envy"
	"log"
)

func main() {
	b := broker.NewBroker(&models.BrokerConfig{
		Host:      envy.Get("HOST", ""),
		Port:      envy.Get("PORT", ""),
		Transport: models.DefaultTransport,
		Storage:   models.InMemoryStorageType,
	})

	//run server
	log.Fatalln(b.Start())

}
