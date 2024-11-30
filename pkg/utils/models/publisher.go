package models

// PublisherConfig sets up connection for the publisher service
type PublisherConfig struct {
	Host      string
	Port      string
	Partition Partition
	Transport Transport
}
