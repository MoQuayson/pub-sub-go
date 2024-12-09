package disk

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"log"
	"os"
	"sync"
)

type LogReader struct {
	mutex       sync.Mutex
	basePath    string
	writer      *bufio.Writer
	file        *os.File
	currentDate string
}

func NewLogReader(basePath string) *LogReader {
	return &LogReader{
		basePath: basePath,
	}
}

func (lr *LogReader) ReadMessages() ([]*models.Message, error) {
	lr.mutex.Lock()
	defer lr.mutex.Unlock()

	messages := make([]*models.Message, 0)
	//read all files from directory
	files, err := os.ReadDir(lr.basePath)
	if err != nil {
		log.Printf("failed to read messages from disk: %v\n", err)
		return nil, err
	}

	for _, file := range files {
		m, err := lr.getMessage(fmt.Sprintf("%s/%s", lr.basePath, file.Name()))
		if err != nil {
			continue
		}

		messages = append(messages, m...)
	}

	return messages, err
}

func (lr *LogReader) getMessage(path string) ([]*models.Message, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	messages := make([]*models.Message, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var msg models.Message
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}
