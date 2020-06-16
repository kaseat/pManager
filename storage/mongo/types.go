package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// Config represents database configuration
type Config struct {
	MongoURL string `json:"mongoURL"`
	DbName   string `json:"dbName"`
}

type dbContext func() context.Context

// Db represents storage
type Db struct {
	syncs      *mongo.Collection
	operations *mongo.Collection
	portfolios *mongo.Collection
	users      *mongo.Collection
	passwords  *mongo.Collection
	context    dbContext
}
