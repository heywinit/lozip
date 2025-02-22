package utils

import (
	"os"
	"strconv"
	"sync"
)

var (
	defaultGenerator *Snowflake
	once            sync.Once
)

// initializeDefaultGenerator initializes the default Snowflake generator
func initializeDefaultGenerator() {
	// Get worker ID and process ID from environment variables or use defaults
	workerID, _ := strconv.ParseInt(os.Getenv("WORKER_ID"), 10, 64)
	processID, _ := strconv.ParseInt(os.Getenv("PROCESS_ID"), 10, 64)

	// Use default values if not set
	if workerID == 0 {
		workerID = 1
	}
	if processID == 0 {
		processID = 1
	}

	var err error
	defaultGenerator, err = NewSnowflake(workerID, processID)
	if err != nil {
		panic(err)
	}
}

// GenerateID generates a new unique ID using the default generator
func GenerateID() int64 {
	once.Do(initializeDefaultGenerator)
	return defaultGenerator.GenerateID()
}

// GenerateStringID generates a string representation of a unique ID
func GenerateStringID() string {
	once.Do(initializeDefaultGenerator)
	return defaultGenerator.GenerateStringID()
} 