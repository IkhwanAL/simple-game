package world

import (
	"log"
	"os"
)

var logger *log.Logger

func InitLogger() {
	file, err := os.OpenFile("world_debug.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("failed to open debug log file: %v", err)
	}

	logger = log.New(file, "", log.LstdFlags)
}

func Logf(format string, args ...any) {
	if logger != nil {
		logger.Printf(format, args...)
	}
}
