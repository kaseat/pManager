package api

type responseStatus string

const (
	ok    responseStatus = "ok"
	notOk responseStatus = "error"
)

type errorResponse struct {
	Status responseStatus `json:"status"`
	Error  string         `json:"error"`
}

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
