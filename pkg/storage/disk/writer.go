package disk

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/MoQuayson/pub-sub-go/pkg/utils"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	fileNameFormat = ""
)

type LogWriter struct {
	mu          sync.Mutex
	basePath    string
	writer      *bufio.Writer
	file        *os.File
	currentDate string
}

func NewLogWriter(basePath string) *LogWriter {
	//check if
	utils.CreateDirectory(basePath)
	return &LogWriter{
		basePath: basePath,
	}
}

func (lw *LogWriter) getLogFilePath() string {
	date := time.Now().Format("2006-01-02")
	return filepath.Join(lw.basePath, fmt.Sprintf("logs-%s.log", date))
}

func (lw *LogWriter) rotateLogIfNeeded() error {
	date := time.Now().Format("2006-01-02")
	if lw.currentDate == date {
		return nil
	}

	// Close the old file
	if lw.file != nil {
		lw.writer.Flush()
		lw.file.Close()
	}

	// Open a new file
	filePath := lw.getLogFilePath()
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	lw.file = file
	lw.writer = bufio.NewWriter(file)
	lw.currentDate = date
	return nil
}

func (lw *LogWriter) WriteMessage(msg *models.Message) error {
	lw.mu.Lock()
	defer lw.mu.Unlock()

	if err := lw.rotateLogIfNeeded(); err != nil {
		return err
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = lw.writer.Write(append(data, '\n'))
	if err != nil {
		return err
	}
	return lw.writer.Flush()
}

func (lw *LogWriter) Close() error {
	lw.mu.Lock()
	defer lw.mu.Unlock()
	if lw.file != nil {
		lw.writer.Flush()
		return lw.file.Close()
	}
	return nil
}
