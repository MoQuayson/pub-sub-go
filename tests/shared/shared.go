package shared

import (
	"github.com/MoQuayson/pub-sub-go/pkg/shared/models"
	"github.com/MoQuayson/pub-sub-go/pkg/shared/utils/enums"
)

var (
	PubConfig = &models.PublisherConfig{
		Host: "0.0.0.0",
		Port: "7000",
	}

	SubConfig = &models.SubscriberConfig{
		Host: "0.0.0.0",
		Port: "7000",
	}

	BrokerCfg = &models.BrokerConfig{
		Host:    "0.0.0.0",
		Port:    "7000",
		Storage: enums.StorageType_InMemory,
	}
)
