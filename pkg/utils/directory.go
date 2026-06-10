package utils

import "os"

func CreateDirectory(path string) error {
	return os.MkdirAll(path, 0750)
}
