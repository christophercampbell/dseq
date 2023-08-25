package app

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"
)

var (
	SinkFilename = "dseq.txt"
)

type Sink struct {
	file *os.File
	path string
}

func NewSink(homeDir string) *Sink {
	path := filepath.Join(homeDir, SinkFilename)
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalf("failed to open sink file %v: %v", path, err)
	}
	return &Sink{file: f, path: path}
}

func (s *Sink) Write(data []byte) {
	hash := common.BytesToHash(data)
	line := fmt.Sprintf("%s\n", hash.Hex())
	if _, err := s.file.WriteString(line); err != nil {
		log.Printf("failed to write data: %v", err)
	}
}

func (s *Sink) Close() {
	err := s.file.Close()
	if err != nil {
		log.Printf("failed to close skin file %v: %v", s.file.Name(), err)
	}
}
