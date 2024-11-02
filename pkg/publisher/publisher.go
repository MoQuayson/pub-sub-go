package publisher

import "github.com/MoQuayson/pub-sub-go/pkg/shared/models"

type Publisher interface {
	PublishMessage(topic string, partition models.Partition, data string) error
}
