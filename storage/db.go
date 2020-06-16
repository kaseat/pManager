package storage

import (
	"time"

	"github.com/kaseat/pManager/models"
	"github.com/kaseat/pManager/storage/mongo"
)

// Db represents data storage
type Db interface {
	SavePassword(user string, hash string) error
	GetPassword(user string) (string, error)
	DeletePassword(user string) (bool, error)
	SaveLastUpdateTime(provider string, date time.Time) error
	GetLastUpdateTime(provider string) (time.Time, error)
	SaveSingleOperation(portfolioID string, op models.Operation) error
	SaveMultipleOperations(portfolioID string, ops []models.Operation) error
	GetOperations(portfolioID string, key string, value string, from string, to string) ([]models.Operation, error)
}

var db mongo.Db

// GetStorage gets storage
func GetStorage() Db {
	if !db.IsInitialized() {
		db = mongo.Db{}
		db.Init(mongo.Config{
			MongoURL: "mongodb://localhost:27017",
			DbName:   "pm_test",
		})
	}
	return db
}
