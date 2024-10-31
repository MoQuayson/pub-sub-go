package models

import (
	"github.com/MoQuayson/go-event-bridge/pkg/shared/utils/enums"
)

type BrokerConfig struct {
	Host    string
	Port    string
	Storage enums.StorageType
}

type RpcConnConfig struct {
	Host, Port string
}
