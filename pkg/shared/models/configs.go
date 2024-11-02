package models

import (
	"github.com/MoQuayson/pub-sub-go/pkg/shared/utils/enums"
)

type BrokerConfig struct {
	Host    string
	Port    string
	Storage enums.StorageType
}

type RpcConnConfig struct {
	Host, Port string
}
