package broker

import (
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	linq "github.com/samber/lo"
	"sort"
	"sync"
	"time"
)

func GetEarliestMessages(mutex *sync.Mutex, subscriberOffsets map[string]models.Offset, messages []*models.Message, request *models.GetMessageRequest) (models.MessageList, error) {
	mutex.Lock()
	offset, exists := subscriberOffsets[request.SubscriberId]
	mutex.Unlock()

	if exists {
		msgCount := len(messages)
		//sort by earliest
		messages = linq.Filter(messages, func(msg *models.Message, index int) bool {
			return msg.Timestamp.Before(offset.Timestamp) || msg.Timestamp.Equal(offset.Timestamp)
		})

		if len(messages) == msgCount {
			messages = nil
		}
	} else {
		//sort by earliest
		sort.Slice(messages, func(i, j int) bool {
			return messages[i].Timestamp.Before(messages[j].Timestamp) || messages[i].Timestamp.Equal(messages[j].Timestamp)
		})
	}

	if len(messages) > 0 {
		mutex.Lock()

		subscriberOffsets[request.SubscriberId] = models.Offset{
			Topic:     request.Topic,
			Partition: request.Partition,
			Timestamp: messages[len(messages)-1].Timestamp,
		}
		mutex.Unlock()
	} else {
		messages = []*models.Message{}
	}

	return messages, nil
}

func GetLatestMessages(mutex *sync.Mutex, subscriberOffsets map[string]models.Offset, messages []*models.Message, request *models.GetMessageRequest) (models.MessageList, error) {
	mutex.Lock()
	offset, exists := subscriberOffsets[request.SubscriberId]
	mutex.Unlock()

	if exists {
		//sort by latest
		messages = linq.Filter(messages, func(msg *models.Message, index int) bool {
			return msg.Timestamp.After(offset.Timestamp)
		})
	} else {
		//sort by latest
		sort.Slice(messages, func(i, j int) bool {
			return messages[i].Timestamp.After(messages[j].Timestamp)
		})
	}
	if len(messages) > 0 {
		mutex.Lock()

		subscriberOffsets[request.SubscriberId] = models.Offset{
			Topic:     request.Topic,
			Partition: request.Partition,
			Timestamp: messages[len(messages)-1].Timestamp,
		}
		mutex.Unlock()
	} else {
		messages = []*models.Message{}
	}

	return messages, nil
}

func GetMessages(mutex *sync.Mutex, subscriberOffsets map[string]models.Offset, messages []*models.Message, request *models.GetMessageRequest, duration time.Duration) (models.MessageList, error) {
	mutex.Lock()
	offset, exists := subscriberOffsets[request.SubscriberId]
	mutex.Unlock()

	if duration < 0 {
		messages = linq.Filter(messages, func(msg *models.Message, index int) bool {
			return msg.Timestamp.Sub(time.Now().UTC().Add(duration)) >= 0
		})
	} else {
		messages = linq.Filter(messages, func(msg *models.Message, index int) bool {

			return time.Now().Sub(msg.Timestamp.Add(duration)) <= 0
		})
	}

	if exists && duration < 0 {
		messages = linq.Filter(messages, func(msg *models.Message, index int) bool {
			return msg.Timestamp.After(offset.Timestamp)
		})
	} else if exists && duration > 0 {
		messages = linq.Filter(messages, func(msg *models.Message, index int) bool {
			return msg.Timestamp.Before(offset.Timestamp)
		})
	}
	if len(messages) > 0 {
		mutex.Lock()

		subscriberOffsets[request.SubscriberId] = models.Offset{
			Topic:     request.Topic,
			Partition: request.Partition,
			Timestamp: messages[len(messages)-1].Timestamp,
		}
		mutex.Unlock()
	} else {
		messages = []*models.Message{}
	}

	return messages, nil
}
