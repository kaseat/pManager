package models

import "time"

// Operation represents market operation
type Operation struct {
	PortfolioID   string        `json:"pid,omitempty"`
	OperationID   string        `json:"id" example:"5edbc0a72c857652a0542fab"`
	Currency      Currency      `json:"currency" example:"USD"`
	Price         float64       `json:"price" example:"293.61"`
	Volume        int64         `json:"vol" example:"100"`
	FIGI          string        `json:"figi" example:"BBG00MVRXDB0"`
	ISIN          string        `json:"isin" example:"BBG00MVRXDB0"`
	Ticker        string        `json:"ticker" example:"VOO"`
	DateTime      time.Time     `json:"date" example:"2020-06-06T15:54:05Z"`
	OperationType OperationType `json:"type" example:"sell"`
}

// Portfolio represets a range of investments
type Portfolio struct {
	PortfolioID string `json:"id" example:"5edb2a0e550dfc5f16392838"`
	UserID      string `json:"userId" example:"5e691429a9bfccacfed4ae2a"`
	Name        string `json:"name" example:"Best portfolio"`
	Description string `json:"description" example:"Best portfolio ever!!!"`
}

// User represents user
type User struct {
	UserID string `json:"id" example:"5edb2a0e550dfc5f16392838"`
	Login  string `json:"login" example:"mark123"`
	Email  string `json:"email" example:"mark123@abc.com"`
}

// Currency represents string currency
type Currency string

// OperationType is market operation type
type OperationType string

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
	// BrokerageFee operation
	BrokerageFee OperationType = "brokerageFee"
	// ExchangeFee operation
	ExchangeFee OperationType = "exchangeFee"
	// PayIn operation
	PayIn OperationType = "payIn"
	// PayOut operation
	PayOut OperationType = "payOut"
	// Coupon operation
	Coupon OperationType = "coupon"
)
