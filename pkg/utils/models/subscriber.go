package models

// SubscriberConfig sets up connection for the subscriber service
type SubscriberConfig struct {
	Host               string
	Port               string
	GroupId            string
	SubscriberId       *string
	MessagePublishTime PublishTime
	Partition          Partition
	Transport          Transport
}
