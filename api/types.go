package api

import (
	"github.com/kaseat/pManager/portfolio"
)

type responseStatus string

const (
	ok    responseStatus = "ok"
	notOk responseStatus = "error"
)

type errorResponse struct {
	Status responseStatus `json:"status"`
	Error  string         `json:"error"`
}

type singlePortfolioIDResponse struct {
	Status      responseStatus `json:"status"`
	PortfolioID string         `json:"createdPortfolioId"`
}

type singlePortfolioUpdateResponse struct {
	Status     responseStatus `json:"status"`
	HasUptated bool           `json:"hasModified"`
}

type portfolioDeleteResponse struct {
	Status     responseStatus `json:"status"`
	HasDeleted bool           `json:"hasDeleted"`
}

type multiplePortfoliosResponse struct {
	Status     responseStatus        `json:"status"`
	Portfolios []portfolio.Portfolio `json:"portfolios"`
}

type singlePortfolioResponse struct {
	Status    responseStatus      `json:"status"`
	Portfolio portfolio.Portfolio `json:"portfolio"`
}

type operationsResponse struct {
	Status      responseStatus        `json:"status"`
	PortfolioID string                `json:"portfolioId"`
	Operations  []portfolio.Operation `json:"operations"`
}

type singleOperationIDResponse struct {
	Status      responseStatus `json:"status"`
	OperationID string         `json:"createdOperationId"`
}

type operationDeleteResponse struct {
	Status     responseStatus `json:"status"`
	NumDeleted int64          `json:"numDeleted"`
}

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
