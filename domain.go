package tcssync

import (
	"time"
)

// Currency represents string currency
type Currency string

// OperationType is market operation type
type OperationType string

// Operation represents market operation
type Operation struct {
	Currency      Currency      `json:"currency"`
	Price         float64       `json:"price"`
	Quantity      int           `json:"quantity"`
	FIGI          string        `json:"figi"`
	DateTime      time.Time     `json:"date"`
	OperationType OperationType `json:"operationType"`
}

// Price represents market instrument price at given time
type Price struct {
	FIGI     string    `json:"figi"`
	Price    float64   `json:"price"`
	Volume   float64   `json:"vol"`
	DateTime time.Time `json:"datetime"`
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
	Buy OperationType = "buy"
	// Sell operation
	Sell OperationType = "sell"
	// BrokerCommission operation
	BrokerCommission OperationType = "brokerCommission"
	// ExchangeCommission operation
	ExchangeCommission OperationType = "exchangeCommission"
	// PayIn operation
	PayIn OperationType = "payIn"
	// PayOut operation
	PayOut OperationType = "payOut"
	// Coupon operation
	Coupon OperationType = "coupon"
)
