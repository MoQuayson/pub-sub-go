package storage

import (
	"github.com/MoQuayson/pub-sub-go/pkg/storage/disk"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"testing"
	"time"
)

func newDiskStorage(t *testing.T) *disk.DiskStorage {
	t.Helper()
	dir := t.TempDir()
	w, err := disk.NewLogWriter(dir)
	if err != nil {
		t.Fatalf("NewLogWriter: %v", err)
	}
	return disk.NewDiskStorage(w, disk.NewLogReader(dir))
}

func TestDiskStorage(t *testing.T) {
	t.Run("store and get roundtrip", func(t *testing.T) {
		s := newDiskStorage(t)
		msg := &models.Message{
			Id:        "MSG001",
			Topic:     "test",
			Partition: models.DefaultPartition,
			Data:      "hello disk",
			Timestamp: time.Now(),
		}

		if err := s.StoreMessage(msg); err != nil {
			t.Fatalf("StoreMessage: %v", err)
		}
		msgs, err := s.GetMessages("test", models.DefaultPartition)
		if err != nil {
			t.Fatalf("GetMessages: %v", err)
		}
		if len(msgs) != 1 || msgs[0].Data != "hello disk" {
			t.Errorf("unexpected messages: %v", msgs)
		}
	})

	t.Run("topic filter", func(t *testing.T) {
		s := newDiskStorage(t)
		_ = s.StoreMessage(&models.Message{Id: "1", Topic: "topicA", Partition: 0, Data: "a", Timestamp: time.Now()})
		_ = s.StoreMessage(&models.Message{Id: "2", Topic: "topicB", Partition: 0, Data: "b", Timestamp: time.Now()})

		msgs, _ := s.GetMessages("topicA", 0)
		if len(msgs) != 1 || msgs[0].Data != "a" {
			t.Errorf("topic filter failed: %v", msgs)
		}
	})

	t.Run("partition filter", func(t *testing.T) {
		s := newDiskStorage(t)
		_ = s.StoreMessage(&models.Message{Id: "1", Topic: "t", Partition: 0, Data: "p0", Timestamp: time.Now()})
		_ = s.StoreMessage(&models.Message{Id: "2", Topic: "t", Partition: 1, Data: "p1", Timestamp: time.Now()})

		msgs, _ := s.GetMessages("t", 0)
		if len(msgs) != 1 || msgs[0].Data != "p0" {
			t.Errorf("partition filter failed: %v", msgs)
		}
	})

	t.Run("EvictMessages returns error not panic", func(t *testing.T) {
		s := newDiskStorage(t)
		err := s.EvictMessages("t", 0, time.Hour)
		if err == nil {
			t.Error("expected error from unimplemented EvictMessages")
		}
	})
}
