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
	OperationID   string    `json:"id" bson:"_id,omitempty"`
	Currency      Currency  `json:"currency"`
	Price         int64     `json:"-" bson:"price"`
	PriceF        float64   `json:"price" bson:"-"`
	Volume        int64     `json:"vol"`
	FIGI          string    `json:"figi"`
	DateTime      time.Time `json:"date"`
	OperationType Type      `json:"operationType"`
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
	PortfolioID string `json:"id" bson:"_id,omitempty"`
	OwnerID     string `json:"ownerId" bson:"oid"`
	Name        string `json:"name"`
	Description string `json:"description"`
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
