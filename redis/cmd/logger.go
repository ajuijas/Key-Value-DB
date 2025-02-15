package cmd

import (
	"log"
	"os"
)

type Logger struct {
	file *os.File
	log *log.Logger
}

func (l *Logger) Close() {
	err := l.file.Sync()
	if err!=nil {
		log.Fatal(err)
	}
	err = l.file.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func NewLogger() *Logger {
	file, err := os.OpenFile(config.StorageFile+"dump.rdb", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log := log.New(file, "", 0)
	return &Logger{file: file, log: log}
}