package storage

import (
	"time"

	"github.com/kaseat/pManager/models"
)

// Db represents data storage
type Db interface {
	SaveLastUpdateTime(provider string, date time.Time) error
	GetLastUpdateTime(provider string) (time.Time, error)
	SaveSingleOperation(portfolioID string, op models.Operation) error
	SaveMultipleOperations(portfolioID string, ops []models.Operation) error
	GetOperations(portfolioID string, key string, value string, from string, to string) ([]models.Operation, error)
}
