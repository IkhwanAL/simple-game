package world

import (
	"encoding/json"
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

func ToJSONBytes(data any) []byte {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Printf("failed to marshal JSON: %v", err)
		return nil
	}
	return jsonBytes
}

func StoreQuickLog(filename string, structInBytes []byte) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("failed to open debug log file: %v", err)
	}

	defer file.Close()

	_, err = file.Write(structInBytes)
	if err != nil {
		log.Fatalf("failed to store snapshot %v", err)
	}
}
