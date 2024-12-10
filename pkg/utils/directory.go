package utils

import "os"

func CreateDirectory(path string) error {
	_ = os.MkdirAll(path, 0750)
	return nil
}
