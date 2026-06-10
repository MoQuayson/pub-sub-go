package models

// SubscriberConfig sets up connection for the subscriber service
type SubscriberConfig struct {
	Host               string
	Port               string
	GroupId            string // reserved for future consumer-group semantics
	SubscriberId       *string
	MessagePublishTime PublishTime
	Partition          Partition
	Transport          Transport
}
