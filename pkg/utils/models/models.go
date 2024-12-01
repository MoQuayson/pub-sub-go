package models

import (
	"time"
)

type Partition int

var (
	DefaultPartition Partition = 0
)

// Message represents the data structure of a message.
type Message struct {
	Id        string
	Topic     string
	Partition Partition
	Data      string
	Timestamp time.Time
}

type MessageList []*Message

type Offset struct {
	Topic     string
	Partition Partition
	Timestamp time.Time
}

type SubscriberOffsets map[string]Offset

type GetMessageRequest struct {
	SubscriberId string
	Topic        string
	Partition    Partition
	Timestamp    time.Time
	PublishTime  PublishTime
}

type PublishTime time.Duration

const (
	LatestPublishTime        PublishTime = 0
	EarliestPublishTime      PublishTime = -1
	WithinASecondPublishTime PublishTime = PublishTime(time.Second * -1)
	WithinAnHourPublishTime  PublishTime = PublishTime(time.Hour * -1)
	WithinADayPublishTime    PublishTime = PublishTime(time.Hour * -24)
)
