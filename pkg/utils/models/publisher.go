package models

// PublisherConfig sets up connection for the publisher service
type PublisherConfig struct {
	Host      string
	Port      string
	Partition Partition // default partition when none is specified at call time
	Transport Transport
}
