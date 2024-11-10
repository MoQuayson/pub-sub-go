package models

import (
	"github.com/MoQuayson/pub-sub-go/pkg/shared/utils/enums"
)

// BrokerConfig sets up connection for the broker service
type BrokerConfig struct {
	Host    string
	Port    string
	Storage enums.StorageType
}
