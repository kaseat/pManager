package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/kaseat/pManager/portfolio"
)

// CreateOwner creates owner
func CreateOwner(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var o portfolio.Owner
	err = json.Unmarshal(body, &o)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if o.Login == "" {
		writeError(w, http.StatusBadRequest, "You must provide login for user")
		return
	}

	co, err := portfolio.AddOwner(o.Login, o.FirstName, o.LastName)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeOk(w, struct {
		Status responseStatus `json:"status"`
		UserID string         `json:"userId"`
	}{
		Status: ok,
		UserID: co.OwnerID,
	})
}
