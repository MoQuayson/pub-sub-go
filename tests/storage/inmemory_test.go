package storage

import (
	"github.com/MoQuayson/pub-sub-go/pkg/storage"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"testing"
	"time"
)

func newMsg(topic string, partition models.Partition, data string, ts time.Time) *models.Message {
	return &models.Message{Id: data, Topic: topic, Partition: partition, Data: data, Timestamp: ts}
}

func TestInMemoryStorage(t *testing.T) {
	t.Run("store and get", func(t *testing.T) {
		s := storage.NewInMemoryStorage()
		msg := newMsg("topic1", 0, "hello", time.Now())

		if err := s.StoreMessage(msg); err != nil {
			t.Fatalf("StoreMessage: %v", err)
		}
		msgs, err := s.GetMessages("topic1", 0)
		if err != nil {
			t.Fatalf("GetMessages: %v", err)
		}
		if len(msgs) != 1 || msgs[0].Data != "hello" {
			t.Errorf("unexpected messages: %v", msgs)
		}
	})

	t.Run("topic isolation", func(t *testing.T) {
		s := storage.NewInMemoryStorage()
		_ = s.StoreMessage(newMsg("topicA", 0, "a", time.Now()))
		_ = s.StoreMessage(newMsg("topicB", 0, "b", time.Now()))

		msgs, _ := s.GetMessages("topicA", 0)
		if len(msgs) != 1 || msgs[0].Data != "a" {
			t.Errorf("topic isolation failed: got %v", msgs)
		}
	})

	t.Run("partition isolation", func(t *testing.T) {
		s := storage.NewInMemoryStorage()
		_ = s.StoreMessage(newMsg("topic", 0, "p0", time.Now()))
		_ = s.StoreMessage(newMsg("topic", 1, "p1", time.Now()))

		msgs, _ := s.GetMessages("topic", 0)
		if len(msgs) != 1 || msgs[0].Data != "p0" {
			t.Errorf("partition isolation failed: got %v", msgs)
		}
	})

	t.Run("evict old messages by TTL", func(t *testing.T) {
		s := storage.NewInMemoryStorage()
		old := newMsg("t", 0, "old", time.Now().Add(-2*time.Hour))
		recent := newMsg("t", 0, "recent", time.Now().Add(-30*time.Minute))
		_ = s.StoreMessage(old)
		_ = s.StoreMessage(recent)

		if err := s.EvictMessages("t", 0, time.Hour); err != nil {
			t.Fatalf("EvictMessages: %v", err)
		}
		msgs, _ := s.GetMessages("t", 0)
		if len(msgs) != 1 || msgs[0].Data != "recent" {
			t.Errorf("eviction wrong: got %v", msgs)
		}
	})

	t.Run("get from empty topic returns nil", func(t *testing.T) {
		s := storage.NewInMemoryStorage()
		msgs, err := s.GetMessages("nonexistent", 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(msgs) != 0 {
			t.Errorf("expected empty, got %d messages", len(msgs))
		}
	})
}
