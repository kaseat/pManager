package models

import (
	"time"

	"github.com/kaseat/pManager/models/currency"
	"github.com/kaseat/pManager/models/instrument"
	"github.com/kaseat/pManager/models/operation"
)

// Price represents price element
type Price struct {
	Price  float64   `json:"price" example:"293.61"`
	Volume int       `json:"vol" example:"100"`
	Date   time.Time `json:"time" example:"2020-06-06T15:54:05Z"`
	ISIN   string    `json:"isin" example:"US9229083632"`
}

// Instrument represents market instrument
type Instrument struct {
	FIGI     string          `json:"figi" example:"BBG000HLJ7M4"`
	ISIN     string          `json:"isin" example:"US45867G1013"`
	Ticker   string          `json:"ticker" example:"IDCC"`
	Name     string          `json:"name" example:"InterDigItal Inc"`
	Type     instrument.Type `json:"type" example:"Stock"`
	Currency currency.Type   `json:"currency" example:"USD"`
}

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

// OperationSorter sorts operations by time.
type OperationSorter []Operation

func (a OperationSorter) Len() int           { return len(a) }
func (a OperationSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a OperationSorter) Less(i, j int) bool { return a[i].DateTime.Unix() < a[j].DateTime.Unix() }
