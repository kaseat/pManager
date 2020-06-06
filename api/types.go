package api

type responseStatus string

const (
	ok    responseStatus = "ok"
	notOk responseStatus = "error"
)

type addPortfoliioSuccess struct {
	PortfolioID string `json:"createdPortfolioId" example:"5edb2a0e550dfc5f16392838"`
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
