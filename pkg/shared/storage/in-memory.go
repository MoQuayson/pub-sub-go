package storage

import (
	"container/list"
	"github.com/MoQuayson/go-event-bridge/pkg/shared/models"
	"sync"
	"time"
)

type InMemoryStorage struct {
	Storage
	messages map[string]map[models.Partition]*list.List
	mutex    sync.Mutex
}

func NewInMemoryStorage() Storage {
	return &InMemoryStorage{messages: make(map[string]map[models.Partition]*list.List)}
}

func (s *InMemoryStorage) StoreMessage(msg *models.Message) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.messages[msg.Topic]; !exists {
		s.messages[msg.Topic] = make(map[models.Partition]*list.List)
	}
	if _, exists := s.messages[msg.Topic][msg.Partition]; !exists {
		s.messages[msg.Topic][msg.Partition] = list.New()
	}

	s.messages[msg.Topic][msg.Partition].PushBack(msg)
	return nil
}

func (s *InMemoryStorage) GetMessages(topic string, partition models.Partition) ([]*models.Message, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if partitionList, exists := s.messages[topic][partition]; exists {
		var msgs []*models.Message
		for e := partitionList.Front(); e != nil; e = e.Next() {
			msgs = append(msgs, e.Value.(*models.Message))
		}
		return msgs, nil
	}
	return nil, nil
}

func (s *InMemoryStorage) EvictMessages(topic string, partition models.Partition, ttl time.Duration) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	partitionList, exists := s.messages[topic][partition]
	if !exists {
		return nil
	}

	now := time.Now()
	for e := partitionList.Front(); e != nil; {
		msg := e.Value.(*models.Message)
		if now.Sub(msg.Timestamp) < ttl {
			break
		}
		next := e.Next()
		partitionList.Remove(e)
		e = next
	}
	return nil
}
