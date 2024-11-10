package models

import (
	"github.com/MoQuayson/pub-sub-go/pkg/shared/utils/enums"
)

type DeliveryResult struct {
	Status enums.DeliveryStatus
	Error  error
}

func NewDeliveryResult(status enums.DeliveryStatus, err error) *DeliveryResult {
	return &DeliveryResult{status, err}
}
