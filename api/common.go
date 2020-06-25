package api

import (
	"encoding/json"
	"net/http"

	"github.com/kaseat/pManager/storage"
)

func writeError(w http.ResponseWriter, statusCode int, text string) {
	resp := errorResponse{Error: text}
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

func canAccess(s storage.Db, login string, pid string) (bool, error) {
	u, err := s.GetUserByLogin(login)
	if err != nil {
		return false, err
	}

	ps, err := s.GetPortfolios(u.UserID)
	if err != nil {
		return false, err
	}

	canAccess := false
	for _, p := range ps {
		if p.PortfolioID == pid {
			canAccess = true
			break
		}
	}
	return canAccess, nil
}
