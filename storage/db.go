package storage

import (
	"time"

	"github.com/kaseat/pManager/models"
	"github.com/kaseat/pManager/models/provider"
	"github.com/kaseat/pManager/storage/mongo"
	"github.com/kaseat/pManager/storage/postgres"
	"golang.org/x/oauth2"
)

// Type - storage type
type Type string

const (
	// Mongo - MongoDB storage
	Mongo Type = "mongo"
	// Postgres - PostgreSQL storage
	Postgres Type = "postgres"
)

// Db represents data storage
type Db interface {
	AddUser(login, email, hash string) (string, error)
	GetUserByLogin(login string) (models.User, error)
	AddUserState(login string, state string) error
	GetUserState(login string) (string, error)
	AddUserToken(state string, token *oauth2.Token) error
	GetUserToken(login string) (oauth2.Token, error)
	GetUserPassword(login string) (string, error)
	UpdateUserPassword(login, hash string) (bool, error)
	DeleteUser(login string) (bool, error)

	AddPortfolio(userID string, p models.Portfolio) (string, error)
	GetPortfolio(userID string, portfolioID string) (models.Portfolio, error)
	GetPortfolios(userID string) ([]models.Portfolio, error)
	UpdatePortfolio(userID string, portfolioID string, p models.Portfolio) (bool, error)
	DeletePortfolio(userID string, portfolioID string) (bool, error)
	DeletePortfolios(userID string) (int64, error)

	AddUserLastUpdateTime(login string, provider provider.Type, date time.Time) error
	GetUserLastUpdateTime(login string, provider provider.Type) (time.Time, error)
	DeleteUserLastUpdateTime(login string, provider provider.Type) error

	AddOperation(portfolioID string, op models.Operation) (string, error)
	AddOperations(portfolioID string, ops []models.Operation) ([]string, error)
	GetOperations(portfolioID string, key string, value string, from string, to string) ([]models.Operation, error)
	DeleteOperation(portfolioID string, operationID string) (bool, error)
	DeleteOperations(portfolioID string) (int64, error)

	AddInstruments(instr []models.Instrument) error
	SetInstrumentPriceUptdTime(sid int, updTime time.Time) (bool, error)
	ClearInstrumentPriceUptdTime(isin string) (bool, error)
	ClearAllInstrumentPriceUptdTime() (bool, error)
	GetInstruments(key string, value string) ([]models.Instrument, error)
	GetAllInstruments() ([]models.Instrument, error)
	DeleteInstruments(key string, value string) (int64, error)
	DeleteAllInstruments() (int64, error)

	AddPrices(prices []models.Price) error
	GetPrices(key, value, from, to string) ([]models.Price, error)
	GetPricesByIsin(isin, from, to string) ([]models.Price, error)
	DeletePrices(key string, value string) (int64, error)
	DeleteAllPrices() (int64, error)

	GetShares(pid string, onDate string) ([]models.Share, error)

	AddTcsToken(token string) error
	DeleteTcsToken() error
	GetTcsToken() (string, error)
}

var dbMongo mongo.Db
var dbPostgres postgres.Db
var currentStorage Type = Postgres

// SwitchStorage switches storage
func SwitchStorage(t Type) {
	currentStorage = t
}

// GetStorage gets storage
func GetStorage() Db {
	switch currentStorage {
	case Postgres:
		if !dbPostgres.IsInitialized() {
			dbPostgres = postgres.Db{}
			dbPostgres.Init(postgres.Config{
				ConnString: "host=localhost port=5432 dbname=p_manager user=test password=test",
			})
		}
		return dbPostgres
	case Mongo:
		if !dbMongo.IsInitialized() {
			dbMongo = mongo.Db{}
			dbMongo.Init(mongo.Config{
				MongoURL: "mongodb://localhost:27017",
				DbName:   "p_manager",
			})
		}
		return dbMongo
	default:
		return nil
	}
}
