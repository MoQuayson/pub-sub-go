package disk

import (
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	linq "github.com/samber/lo"
	"time"
)

type DiskStorage struct {
	Directory string
	writer    *LogWriter
	reader    *LogReader
}

func NewDiskStorage(w *LogWriter, r *LogReader) *DiskStorage {
	return &DiskStorage{writer: w, reader: r}
}

func (d *DiskStorage) StoreMessage(msg *models.Message) error {
	return d.writer.WriteMessage(msg)
}

func (d *DiskStorage) GetMessages(topic string, partition models.Partition) (models.MessageList, error) {
	//get messages
	messages, err := d.reader.ReadMessages()
	if err != nil {
		return nil, err
	}

	filteredMessages := linq.Filter(messages, func(item *models.Message, index int) bool {
		return item.Topic == topic && item.Partition == partition
	})
	return filteredMessages, nil
}

func (d *DiskStorage) EvictMessages(topic string, partition models.Partition, ttl time.Duration) error {
	//TODO implement me
	panic("implement me")
}
