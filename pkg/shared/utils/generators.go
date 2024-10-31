package utils

import (
	"fmt"
	"github.com/gofrs/uuid"
	"strings"
)

func NewUUID() uuid.UUID {
	uid, _ := uuid.NewV4()

	return uid
}

func NewMessageId() string {
	return fmt.Sprintf("MSG%s", strings.ToUpper(strings.Replace(NewUUID().String(), "-", "", -1)))
}

func NewSubscriberId() string {
	return fmt.Sprintf("SUB%s", strings.ToUpper(strings.Replace(NewUUID().String(), "-", "", -1)))
}
