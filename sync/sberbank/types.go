package sberbank

import (
	"time"

	"github.com/kaseat/pManager/models"
)

type pair struct {
	begin, end int
}

type ticker string

type operationInfo struct {
	Currency      string
	Price         float64
	Volume        int64
	ISIN          string
	Ticker        string
	OperationTime time.Time
	OperationType string
}

type securitiesInfo struct {
	Ticker string
	ISIN   string
	IsBond bool
}

type report struct {
	IsEmpty        bool
	Date           string
	Operations     []models.Operation
	SecuritiesInfo map[ticker]securitiesInfo
	CashFlow       []models.Operation
	Buybacks       []models.Operation
}
