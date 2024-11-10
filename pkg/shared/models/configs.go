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

// SubscriberConfig sets up connection for the subscriber service
type SubscriberConfig struct {
	Host               string
	Port               string
	GroupId            string
	SubscriberId       *string
	MessagePublishTime enums.PublishTime
	Partition          Partition
}

// PublisherConfig sets up connection for the publisher service
type PublisherConfig struct {
	Host      string
	Port      string
	Partition Partition
}
