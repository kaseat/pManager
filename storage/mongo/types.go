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
	syncs       *mongo.Collection
	operations  *mongo.Collection
	portfolios  *mongo.Collection
	users       *mongo.Collection
	instruments *mongo.Collection
	settings    *mongo.Collection
	context     dbContext
}

type token struct {
	AccessToken  string `bson:"access_token"`
	TokenType    string `bson:"token_type,omitempty"`
	RefreshToken string `bson:"refresh_token,omitempty"`
	Expiry       string `bson:"expiry,omitempty"`
}
