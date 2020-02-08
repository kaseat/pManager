package api

import (
	"encoding/json"
	"net/http"
)

func writeError(w http.ResponseWriter, statusCode int, text string) {
	resp := errorResponse{Status: notOk, Error: text}
	bytes, _ := json.Marshal(&resp)
	w.WriteHeader(statusCode)
	w.Write(bytes)
}
