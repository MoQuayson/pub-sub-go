package broker

import (
	"github.com/MoQuayson/pub-sub-go/internal/broker"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"sync"
	"testing"
	"time"
)

func makeMessages(n int, baseTime time.Time) []*models.Message {
	msgs := make([]*models.Message, n)
	for i := range msgs {
		msgs[i] = &models.Message{
			Id:        "MSG",
			Topic:     "test",
			Partition: models.DefaultPartition,
			Data:      "data",
			Timestamp: baseTime.Add(time.Duration(i) * time.Second),
		}
	}
	return msgs
}

func TestGetLatestMessages(t *testing.T) {
	base := time.Now().Truncate(time.Second)

	t.Run("no prior offset returns all messages", func(t *testing.T) {
		msgs := makeMessages(3, base)
		req := &models.GetMessageRequest{SubscriberId: "sub1", Topic: "test", Partition: models.DefaultPartition}
		var mu sync.Mutex
		offsets := models.SubscriberOffsets{}

		got, err := broker.GetLatestMessages(&mu, offsets, msgs, req)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 3 {
			t.Errorf("expected 3 messages, got %d", len(got))
		}
	})

	t.Run("with prior offset returns only newer messages", func(t *testing.T) {
		msgs := makeMessages(3, base)
		req := &models.GetMessageRequest{SubscriberId: "sub2", Topic: "test", Partition: models.DefaultPartition}
		var mu sync.Mutex
		offsets := models.SubscriberOffsets{
			broker.OffsetKey(req): {Topic: "test", Partition: models.DefaultPartition, Timestamp: msgs[1].Timestamp},
		}

		got, err := broker.GetLatestMessages(&mu, offsets, msgs, req)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 1 {
			t.Errorf("expected 1 message after offset, got %d", len(got))
		}
	})

	t.Run("no messages returns empty slice", func(t *testing.T) {
		req := &models.GetMessageRequest{SubscriberId: "sub3", Topic: "test", Partition: models.DefaultPartition}
		var mu sync.Mutex
		offsets := models.SubscriberOffsets{}

		got, err := broker.GetLatestMessages(&mu, offsets, nil, req)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 0 {
			t.Errorf("expected empty slice, got %d messages", len(got))
		}
	})
}

func TestSubscriberOffsetsAreScopedByTopicAndPartition(t *testing.T) {
	base := time.Now().Truncate(time.Second)
	reqPartition0 := &models.GetMessageRequest{SubscriberId: "sub-scope", Topic: "test", Partition: 0}
	reqPartition1 := &models.GetMessageRequest{SubscriberId: "sub-scope", Topic: "test", Partition: 1}
	msgs := makeMessages(3, base)
	var mu sync.Mutex
	offsets := models.SubscriberOffsets{
		broker.OffsetKey(reqPartition0): {Topic: "test", Partition: 0, Timestamp: msgs[1].Timestamp},
	}

	got, err := broker.GetLatestMessages(&mu, offsets, msgs, reqPartition1)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 {
		t.Fatalf("expected independent partition offset to return 3 messages, got %d", len(got))
	}
	if _, exists := offsets[broker.OffsetKey(reqPartition0)]; !exists {
		t.Fatal("expected partition 0 offset to remain")
	}
	if _, exists := offsets[broker.OffsetKey(reqPartition1)]; !exists {
		t.Fatal("expected partition 1 offset to be tracked independently")
	}
}

func TestGetEarliestMessages(t *testing.T) {
	base := time.Now().Truncate(time.Second)

	t.Run("no prior offset returns all messages sorted oldest first", func(t *testing.T) {
		msgs := makeMessages(3, base)
		req := &models.GetMessageRequest{SubscriberId: "sub4", Topic: "test", Partition: models.DefaultPartition}
		var mu sync.Mutex
		offsets := models.SubscriberOffsets{}

		got, err := broker.GetEarliestMessages(&mu, offsets, msgs, req)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 3 {
			t.Errorf("expected 3 messages, got %d", len(got))
		}
		if !got[0].Timestamp.Before(got[len(got)-1].Timestamp) && !got[0].Timestamp.Equal(got[len(got)-1].Timestamp) {
			t.Error("messages not sorted oldest-first")
		}
	})

	t.Run("no messages returns empty slice", func(t *testing.T) {
		req := &models.GetMessageRequest{SubscriberId: "sub5", Topic: "test", Partition: models.DefaultPartition}
		var mu sync.Mutex
		offsets := models.SubscriberOffsets{}

		got, err := broker.GetEarliestMessages(&mu, offsets, nil, req)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 0 {
			t.Errorf("expected empty slice, got %d messages", len(got))
		}
	})
}

func TestGetMessages_WithDuration(t *testing.T) {
	t.Run("negative duration returns messages within window", func(t *testing.T) {
		// Messages: one from 2 hours ago, one from 30 minutes ago
		now := time.Now()
		msgs := []*models.Message{
			{Id: "old", Topic: "test", Partition: models.DefaultPartition, Timestamp: now.Add(-2 * time.Hour)},
			{Id: "recent", Topic: "test", Partition: models.DefaultPartition, Timestamp: now.Add(-30 * time.Minute)},
		}
		req := &models.GetMessageRequest{SubscriberId: "sub6", Topic: "test", Partition: models.DefaultPartition}
		var mu sync.Mutex
		offsets := models.SubscriberOffsets{}

		got, err := broker.GetMessages(&mu, offsets, msgs, req, -1*time.Hour)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 1 || got[0].Id != "recent" {
			t.Errorf("expected 1 recent message, got %d", len(got))
		}
	})
}
