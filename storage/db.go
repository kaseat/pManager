package storage

import (
	"time"

	"github.com/kaseat/pManager/models"
	"github.com/kaseat/pManager/storage/mongo"
)

// Db represents data storage
type Db interface {
	AddUser(login, email, hash string) (string, error)
	GetUserByLogin(login string) (models.User, error)
	GetUserPassword(login string) (string, error)
	UpdateUser(login string, user models.User) (bool, error)
	UpdateUserPassword(login, hash string) (bool, error)
	DeleteUser(login string) (bool, error)

	AddPortfolio(userID string, p models.Portfolio) (string, error)
	GetPortfolio(userID string, portfolioID string) (models.Portfolio, error)
	GetPortfolios(userID string) ([]models.Portfolio, error)
	UpdatePortfolio(userID string, portfolioID string, p models.Portfolio) (bool, error)
	DeletePortfolio(userID string, portfolioID string) (bool, error)
	DeletePortfolios(userID string) (int64, error)

	AddLastUpdateTime(provider string, date time.Time) error
	GetLastUpdateTime(provider string) (time.Time, error)
	DeleteLastUpdateTime(provider string) error

	AddOperation(portfolioID string, op models.Operation) (string, error)
	AddOperations(portfolioID string, ops []models.Operation) ([]string, error)
	GetOperations(portfolioID string, key string, value string, from string, to string) ([]models.Operation, error)
	DeleteOperation(portfolioID string, operationID string) (bool, error)
	DeleteOperations(portfolioID string) (int64, error)
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
