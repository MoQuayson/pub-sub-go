package main

import (
	"github.com/MoQuayson/pub-sub-go/pkg/subscriber"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"github.com/gobuffalo/envy"
	"log"
)

func main() {
	//sub := subscriber.NewSubscriber(&models.SubscriberConfig{
	//	Host:               envy.Get("HOST", ""),
	//	Port:               envy.Get("PORT", ""),
	//	MessagePublishTime: enums.PublishTime_Latest,
	//	Partition:          models.DefaultPartition,
	//})

	//subId := "SUB1234"
	//sub := subscriber.NewSubscriber(&models.SubscriberConfig{
	//	Host:               envy.Get("HOST", ""),
	//	Port:               envy.Get("PORT", ""),
	//	MessagePublishTime: enums.PublishTime_Earliest,
	//	Partition:          models.DefaultPartition,
	//	SubscriberId:       &subId,
	//})

	sub := subscriber.NewSubscriber(&models.SubscriberConfig{
		Host:               envy.Get("HOST", ""),
		Port:               envy.Get("PORT", ""),
		MessagePublishTime: models.EarliestPublishTime,
		Partition:          models.DefaultPartition,
		Transport:          models.GrpcTransport,
		//SubscriberId:       &subId,
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
