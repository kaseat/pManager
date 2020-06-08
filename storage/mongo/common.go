package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Init mongodb module
func (db *Db) Init(config Config) error {
	cfg := config
	db.context = func() context.Context { return context.Background() }
	clientOptions := options.Client().ApplyURI(cfg.MongoURL)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return err
	}
	err = client.Connect(db.context())
	if err != nil {
		return err
	}
	db.syncs = client.Database(cfg.DbName).Collection("syncs")
	db.operations = client.Database(cfg.DbName).Collection("operations")
	db.portfolios = client.Database(cfg.DbName).Collection("portfolios")
	return nil
}
