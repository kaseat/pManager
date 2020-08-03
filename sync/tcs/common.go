package tcs

import (
	"time"
)

// Status represents sync status
type Status string

// SyncStatus represents current snc status
type SyncStatus struct {
	Status Status `json:"status"`
	Error  error  `json:"error,omitempty"`
}

type syncError struct {
	Error      error
	IsNotEmpty bool
}

type chunk struct {
	From time.Time
	To   time.Time
}

const (
	// Ok shows that last sync was ok
	Ok Status = "ok"
	// Err shows that there was error during sync
	Err Status = "Error"
	// Processing shows that sync is in process
	Processing Status = "Processing"
)

func today() time.Time {
	year, month, day := time.Now().Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func getTimeChunks(from, to time.Time, size int) []chunk {
	result := []chunk{}
	for {
		end := from.AddDate(0, size, 0)
		if to.Before(end) {
			result = append(result, chunk{From: from, To: to})
			break
		}
		result = append(result, chunk{From: from, To: end})
		from = end.AddDate(0, 0, 1)
	}
	return result
}
