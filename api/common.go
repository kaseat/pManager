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

func writeOk(w http.ResponseWriter, resp interface{}) {
	bytes, err := json.Marshal(&resp)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
