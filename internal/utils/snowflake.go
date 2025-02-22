package utils

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

const (
	// Discord Epoch (2015-01-01T00:00:00Z)
	epoch int64 = 1420070400000

	// Bits allocated for each component
	timestampBits = 42
	workerIDBits  = 5
	processIDBits = 5
	sequenceBits  = 12

	// Maximum values for each component
	maxWorkerID  = -1 ^ (-1 << workerIDBits)
	maxProcessID = -1 ^ (-1 << processIDBits)
	maxSequence  = -1 ^ (-1 << sequenceBits)

	// Bit shifts for each component
	workerIDShift  = sequenceBits
	processIDShift = sequenceBits + workerIDBits
	timestampShift = sequenceBits + workerIDBits + processIDBits
)

type Snowflake struct {
	mutex     sync.Mutex
	timestamp int64
	workerID  int64
	processID int64
	sequence  int64
}

// NewSnowflake creates a new Snowflake instance
func NewSnowflake(workerID, processID int64) (*Snowflake, error) {
	if workerID < 0 || workerID > maxWorkerID {
		return nil, fmt.Errorf("worker ID must be between 0 and %d", maxWorkerID)
	}
	if processID < 0 || processID > maxProcessID {
		return nil, fmt.Errorf("process ID must be between 0 and %d", maxProcessID)
	}

	return &Snowflake{
		workerID:  workerID,
		processID: processID,
	}, nil
}

// GenerateID generates a new unique Snowflake ID
func (s *Snowflake) GenerateID() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	timestamp := time.Now().UnixMilli() - epoch

	if timestamp == s.timestamp {
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 {
			// Sequence exhausted, wait for next millisecond
			for timestamp <= s.timestamp {
				timestamp = time.Now().UnixMilli() - epoch
			}
		}
	} else {
		s.sequence = 0
	}

	s.timestamp = timestamp

	id := (timestamp << timestampShift) |
		(s.processID << processIDShift) |
		(s.workerID << workerIDShift) |
		s.sequence

	return id
}

// ExtractTimestamp extracts the timestamp from a Snowflake ID
func ExtractTimestamp(id int64) time.Time {
	timestamp := (id >> timestampShift) + epoch
	return time.UnixMilli(timestamp)
}

// GenerateStringID generates a string representation of a Snowflake ID
func (s *Snowflake) GenerateStringID() string {
	return strconv.FormatInt(s.GenerateID(), 10)
} 