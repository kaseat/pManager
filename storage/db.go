package storage

import (
	"time"

	"github.com/kaseat/pManager/models"
)

// Db represents data storage
type Db interface {
	SaveLastUpdateTime(provider string, date time.Time) error
	GetLastUpdateTime(provider string) (time.Time, error)
	SaveMultipleOperations(portfolioID string, ops []models.Operation) error
}
