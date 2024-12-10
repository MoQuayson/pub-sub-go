# PubSubGo

PubSubGo is a high-performance, scalable in-memory Pub/Sub (Publisher/Subscriber) messaging system implemented in Go.
This package provides low-latency message broadcasting to multiple subscribers, with configurable message storage backends including in-memory cache, Redis, and disk storage.

## Features

- **In-Memory Pub/Sub**: Fast and lightweight for high-throughput messaging.
- **Configurable Storage**: Choose between in-memory, Redis, or disk storage based on use case.
- **Scalable Architecture**: Designed to scale horizontally for distributed environments.
- **Flexible Communication**: Expandable options beyond RPC for custom implementations.

## Installation

```bash
go get https://github.com/MoQuayson/pub-sub-go@latest
```

## Getting Started

### Initiate Broker

```go
package main

import (
	"github.com/MoQuayson/pub-sub-go/pkg/broker"
	"github.com/MoQuayson/pub-sub-go/pkg/utils"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"github.com/gobuffalo/envy"
	"log"
)

func main() {
	b := broker.NewBroker(&models.BrokerConfig{
		Host:      envy.Get("HOST", "127.0.0.1"),
		Port:      envy.Get("PORT", "7000"),
		Transport: models.DefaultTransport,
		Storage:   models.InMemoryStorageType,
	})

	log.Fatalln(b.Start())
}
```

### Initialize Publisher

```go
package main

import (
	"github.com/MoQuayson/pub-sub-go/pkg/broker"
	"github.com/MoQuayson/pub-sub-go/pkg/utils"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"github.com/gobuffalo/envy"
	"log"
)

func main() {
	pub := publisher.NewPublisher(&models.PublisherConfig{
		Host:      envy.Get("HOST", ""),
		Port:      envy.Get("PORT", ""),
		Transport: models.DefaultTransport,
	})

	// publish message
	if err := pub.PublishMessage("test", models.DefaultPartition, "hello world"); err != nil {
		log.Fatalln("failed to publish message: ", err)
	}
}
```


### Inititialize Subscriber

```go
package main

import (
	"github.com/MoQuayson/pub-sub-go/pkg/broker"
	"github.com/MoQuayson/pub-sub-go/pkg/utils"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"github.com/gobuffalo/envy"
	"log"
)

func main() {
	sub := subscriber.NewSubscriber(&models.SubscriberConfig{
		Host:               envy.Get("HOST", ""),
		Port:               envy.Get("PORT", ""),
		MessagePublishTime: models.EarliestPublishTime,
		Partition:          models.DefaultPartition,
		Transport:          models.DefaultTransport,
	})
	for {
		messages, err := sub.Subscribe("test")
		if err != nil {
			log.Fatalf("failed to subscribe message: %v\n", err)
		}

		for _, message := range messages {
			log.Printf("Message (%s) with data:%v\n", message.Id, message.Data)
		}
	}
}
```

## Transport Options
* Default (RPC)
* Redis
* Grpc

## Storage Options
* InMemory: Fastest option, ideal for transient messaging.
* Redis: Suitable for distributed setups where persistence and reliability are required.
* Disk: Provides durability for long-term storage.

## Publish Time Options
* Earliest: retrieves messages from the dawn of time
* Latest: gets only current messages produced to a particular topic and partition 
* WithinAnHour: gets data within an hour timeframe
* WithinADay: gets messages within a day's timeframe
* Custom: gets messages based on the duration given. e.g. PublishTime( models.PublishTime(time.Minute * 5)) //get last 5 minutes old messages

## Performance
PubSubGo is optimized for low-latency, high-throughput messaging making it suitable for real-time applications.