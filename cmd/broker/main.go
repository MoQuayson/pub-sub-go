package main

import (
	"github.com/MoQuayson/go-event-bridge/internal/rpc/broker"
	"github.com/MoQuayson/go-event-bridge/pkg/shared/models"
	"github.com/MoQuayson/go-event-bridge/pkg/shared/utils/enums"
	"github.com/gobuffalo/envy"
)

func main() {
	b := broker.NewRpcBroker(&models.BrokerConfig{
		Host:    envy.Get("HOST", ""),
		Port:    envy.Get("PORT", ""),
		Storage: enums.StorageType_InMemory,
	})

	//run server
	b.Start()

}
