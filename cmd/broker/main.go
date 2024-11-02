package main

import (
	"github.com/MoQuayson/pub-sub-go/pkg/broker"
	"github.com/MoQuayson/pub-sub-go/pkg/shared/models"
	"github.com/MoQuayson/pub-sub-go/pkg/shared/utils/enums"
	"github.com/gobuffalo/envy"
)

func main() {
	b := broker.NewBrokerService(&models.BrokerConfig{
		Host:    envy.Get("HOST", ""),
		Port:    envy.Get("PORT", ""),
		Storage: enums.StorageType_InMemory,
	})

	//run server
	b.Start()

}
