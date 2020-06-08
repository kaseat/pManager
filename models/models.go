package models

import "time"

// Operation represents market operation
type Operation struct {
	PortfolioID   string    `json:"pid,omitempty"`
	OperationID   string    `json:"id" example:"5edbc0a72c857652a0542fab"`
	Currency      Currency  `json:"currency" example:"USD"`
	Price         float64   `json:"price" example:"293.61"`
	Volume        int64     `json:"vol" example:"100"`
	FIGI          string    `json:"figi" example:"BBG00MVRXDB0"`
	ISIN          string    `json:"isin" example:"BBG00MVRXDB0"`
	Ticker        string    `json:"ticker" example:"VOO"`
	DateTime      time.Time `json:"date" example:"2020-06-06T15:54:05Z"`
	OperationType Type      `json:"type" example:"sell"`
}

// Currency represents string currency
type Currency string

// Type is market operation type
type Type string

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
	// BrokerageFee operation
	BrokerageFee Type = "brokerageFee"
	// ExchangeFee operation
	ExchangeFee Type = "exchangeFee"
	// PayIn operation
	PayIn Type = "payIn"
	// PayOut operation
	PayOut Type = "payOut"
	// Coupon operation
	Coupon Type = "coupon"
)
