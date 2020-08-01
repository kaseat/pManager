package tcs

// Status represents sync status
type Status string

const (
	// Ok shows that last sync was ok
	Ok Status = "ok"
	// Err shows that there was error during sync
	Err Status = "Error"
	// Processing shows that sync is in process
	Processing Status = "Processing"
)

// SyncStatus represents current snc status
type SyncStatus struct {
	Status Status `json:"status"`
	Error  error  `json:"error,omitempty"`
}
