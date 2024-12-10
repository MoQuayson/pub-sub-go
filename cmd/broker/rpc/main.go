package main

import (
	"github.com/MoQuayson/pub-sub-go/pkg/broker"
	"github.com/MoQuayson/pub-sub-go/pkg/utils"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"github.com/gobuffalo/envy"
	"log"
)

func main() {
	//b:= getBrokerWithInMemoryStorage()
	b := getBrokerWithDiskStorage()

	//run server
	log.Fatalln(b.Start())

}

// in-memory broker example
func getBrokerWithInMemoryStorage() broker.Broker {
	return broker.NewBroker(&models.BrokerConfig{
		Host:      envy.Get("HOST", "127.0.0.1"),
		Port:      envy.Get("PORT", "7000"),
		Transport: models.DefaultTransport,
		Storage:   models.InMemoryStorageType,
	})
}

// disk broker example
func getBrokerWithDiskStorage() broker.Broker {
	return broker.NewBroker(&models.BrokerConfig{
		Host:            envy.Get("HOST", "127.0.0.1"),
		Port:            envy.Get("PORT", "7000"),
		Transport:       models.DefaultTransport,
		Storage:         models.DiskStorageType,
		StorageLocation: utils.ConvertToPointerString("./bin/logs"),
	})
}
