package portfolio

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// Currency represents string currency
type Currency string

// Type is market operation type
type Type string

// Operation represents market operation
type Operation struct {
	PortfolioID   string    `json:"pid,omitempty" bson:"portfolio"`
	OperationID   string    `json:"id" bson:"_id,omitempty" example:"5edbc0a72c857652a0542fab"`
	Currency      Currency  `json:"currency" example:"USD"`
	Price         int64     `json:"-" bson:"price"`
	PriceF        float64   `json:"price" bson:"-" example:"293.61"`
	Volume        int64     `json:"vol" example:"100"`
	FIGI          string    `json:"figi" example:"BBG00MVRXDB0"`
	DateTime      time.Time `json:"date" example:"2020-06-06T15:54:05Z"`
	OperationType Type      `json:"operationType" example:"sell"`
}

// Owner represets owner of a portfolio
type Owner struct {
	OwnerID   string `json:"id" bson:"_id,omitempty"`
	Login     string `json:"login"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// Portfolio represets a range of investments
type Portfolio struct {
	PortfolioID string `json:"id" bson:"_id,omitempty" example:"5edb2a0e550dfc5f16392838"`
	OwnerID     string `json:"ownerId" bson:"oid" example:"5e691429a9bfccacfed4ae2a"`
	Name        string `json:"name" example:"Best portfolio"`
	Description string `json:"description" example:"Best portfolio ever!!!"`
}

// Config represents configuration
type Config struct {
	MongoURL string `json:"mongoURL"`
	DbName   string `json:"dbName"`
}

type dbContext func() context.Context

type database struct {
	operations *mongo.Collection
	prices     *mongo.Collection
	portfolios *mongo.Collection
	owners     *mongo.Collection
	context    dbContext
}

const (
	// RUB is Russian Rubble
	RUB Currency = "RUB"
	// USD is US Dollar
	USD Currency = "USD"
	// EUR is Euro
	EUR Currency = "EUR"
)

const (
	// Buy operation
	Buy Type = "buy"
	// Sell operation
	Sell Type = "sell"
	// BrokerCommission operation
	BrokerCommission Type = "brokerCommission"
	// ExchangeCommission operation
	ExchangeCommission Type = "exchangeCommission"
	// PayIn operation
	PayIn Type = "payIn"
	// PayOut operation
	PayOut Type = "payOut"
	// Coupon operation
	Coupon Type = "coupon"
)
