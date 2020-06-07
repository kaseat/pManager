package api

type responseStatus string

const (
	ok    responseStatus = "ok"
	notOk responseStatus = "error"
)

type addPortfoliioSuccess struct {
	PortfolioID string `json:"createdPortfolioId" example:"5edb2a0e550dfc5f16392838"`
}

type addOperationSuccess struct {
	OperationID string `json:"createdOperationId" example:"5edbc0a72c857652a0542fab"`
}

type putPortfoliioSuccess struct {
	HasModified bool `json:"hasModified" example:"true"`
}

type delPortfoliioSuccess struct {
	HasDeleted bool `json:"hasDeleted" example:"true"`
}

type portfolioRequest struct {
	Name        string `json:"name" example:"Best portfolio"`
	Description string `json:"description" example:"Best portfolio ever!!!"`
}

type operationRequest struct {
	Currency      string  `json:"currency" example:"USD"`
	Price         float64 `json:"price" example:"293.61"`
	Volume        int64   `json:"vol" example:"100"`
	FIGI          string  `json:"figi" example:"BBG00MVRXDB0"`
	DateTime      string  `json:"date" example:"2020-06-06T22:54:05.000+07:00"`
	OperationType string  `json:"operationType" example:"sell"`
}

type errorResponse struct {
	Error string `json:"error" example:"Something went wrong"`
}

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type tokenResponse struct {
	Status responseStatus `json:"status"`
	Token  string         `json:"token"`
}
