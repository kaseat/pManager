package models

import (
	"time"

	"github.com/kaseat/pManager/models/currency"
	"github.com/kaseat/pManager/models/operation"
)

// Operation represents market operation
type Operation struct {
	PortfolioID   string         `json:"pid,omitempty"`
	OperationID   string         `json:"id" example:"5edbc0a72c857652a0542fab"`
	Currency      currency.Type  `json:"currency" example:"USD"`
	Price         float64        `json:"price" example:"293.61"`
	Volume        int64          `json:"vol" example:"100"`
	FIGI          string         `json:"figi,omitempty" example:"BBG00MVRXDB0"`
	ISIN          string         `json:"isin,omitempty" example:"US9229083632"`
	Ticker        string         `json:"ticker" example:"VOO"`
	DateTime      time.Time      `json:"date" example:"2020-06-06T15:54:05Z"`
	OperationType operation.Type `json:"type" example:"sell"`
}

// Portfolio represets a range of investments
type Portfolio struct {
	PortfolioID string `json:"id" example:"5edb2a0e550dfc5f16392838"`
	UserID      string `json:"-"`
	Name        string `json:"name" example:"Best portfolio"`
	Description string `json:"description" example:"Best portfolio ever!!!"`
}

// User represents user
type User struct {
	UserID string `json:"id" example:"5edb2a0e550dfc5f16392838"`
	Login  string `json:"login" example:"mark123"`
	Email  string `json:"email" example:"mark123@abc.com"`
}
