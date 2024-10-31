package models

import (
	"github.com/MoQuayson/go-event-bridge/pkg/shared/utils/enums"
)

type PublishResponse struct {
	Status enums.DeliveryStatus
	Data   any
}

type Empty struct {
}

type StringValue struct {
	Value string
}

func NewPublishResponse(status enums.DeliveryStatus, data any) *PublishResponse {
	return &PublishResponse{Status: status, Data: data}
}
