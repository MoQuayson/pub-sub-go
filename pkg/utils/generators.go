package utils

import (
	"fmt"
	"github.com/gofrs/uuid"
	"strings"
)

func NewUUID() (uuid.UUID, error) {
	return uuid.NewV4()
}

func NewMessageId() (string, error) {
	uid, err := NewUUID()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("MSG%s", strings.ToUpper(strings.Replace(uid.String(), "-", "", -1))), nil
}

func NewPublisherId() (string, error) {
	uid, err := NewUUID()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("PUB%s", strings.ToUpper(strings.Replace(uid.String(), "-", "", -1))), nil
}

func NewSubscriberId() (string, error) {
	uid, err := NewUUID()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("SUB%s", strings.ToUpper(strings.Replace(uid.String(), "-", "", -1))), nil
}
