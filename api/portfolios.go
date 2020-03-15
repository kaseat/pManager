package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kaseat/pManager/portfolio"
)

// CreateSinglePortfolio creates single portfolio
func CreateSinglePortfolio(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	user := r.Header.Get("user")

	found, o, err := portfolio.GetOwnerByLogin(user)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if !found {
		writeError(w, http.StatusNotFound, fmt.Sprint("No user found with login: ", user))
		return
	}
	var p portfolio.Portfolio
	err = json.Unmarshal(body, &p)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if p.Name == "" {
		writeError(w, http.StatusBadRequest, "You must provide portfolio name")
		return
	}

	p, err = o.AddPortfolio(p.Name, p.Description)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeOk(w, struct {
		Status      responseStatus `json:"status"`
		PortfolioID string         `json:"createdPortfolioId"`
	}{Status: ok, PortfolioID: p.PortfolioID})
}

// ReadSinglePortfolio gets single portfolio by id
func ReadSinglePortfolio(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	user := r.Header.Get("user")

	found, o, err := portfolio.GetOwnerByLogin(user)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if !found {
		writeError(w, http.StatusNotFound, fmt.Sprint("No user found with login: ", user))
		return
	}

	found, p, err := o.GetPortfolio(id)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if !found {
		writeError(w, http.StatusNotFound, fmt.Sprint("No portfolio found with Id: ", id))
		return
	}

	writeOk(w, struct {
		Status    responseStatus      `json:"status"`
		Portfolio portfolio.Portfolio `json:"portfolio"`
	}{Status: ok, Portfolio: p})
}

// ReadAllPortfolios gets all portfolios
func ReadAllPortfolios(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user := r.Header.Get("user")

	found, o, err := portfolio.GetOwnerByLogin(user)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if !found {
		writeError(w, http.StatusNotFound, fmt.Sprint("No user found with login: ", user))
		return
	}

	ps, err := o.GetAllPortfolios()
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(ps) == 0 {
		writeError(w, http.StatusNotFound, fmt.Sprint("No portfolios found"))
		return
	}

	writeOk(w, struct {
		Status     responseStatus        `json:"status"`
		Portfolios []portfolio.Portfolio `json:"portfolios"`
	}{Status: ok, Portfolios: ps})
}

// UptateSinglePortfolio updates single portfolio by id
func UptateSinglePortfolio(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var p portfolio.Portfolio

	err := decoder.Decode(&p)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if p.Name == "" {
		writeError(w, http.StatusBadRequest, "You must provide portfolio name")
		return
	}

	p.PortfolioID = mux.Vars(r)["id"]
	hasUptated, err := p.UpdatePortfolio()
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeOk(w, struct {
		Status     responseStatus `json:"status"`
		HasUptated bool           `json:"hasModified"`
	}{Status: ok, HasUptated: hasUptated})
}

// DeleteSinglePortfolio deletes single portfolio by id
func DeleteSinglePortfolio(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	user := r.Header.Get("user")

	found, o, err := portfolio.GetOwnerByLogin(user)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if !found {
		writeError(w, http.StatusNotFound, fmt.Sprint("No user found with login: ", user))
		return
	}

	hasDeleted, err := o.DeletePortfolio(id)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if !hasDeleted {
		writeError(w, http.StatusNotFound, fmt.Sprint("No portfolio found with Id: ", id))
		return
	}

	writeOk(w, struct {
		Status     responseStatus `json:"status"`
		HasDeleted bool           `json:"hasDeleted"`
	}{Status: ok, HasDeleted: true})
}

// DeleteAllPortfolios deletes all portfolios
func DeleteAllPortfolios(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user := r.Header.Get("user")

	found, o, err := portfolio.GetOwnerByLogin(user)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if !found {
		writeError(w, http.StatusNotFound, fmt.Sprint("No user found with login: ", user))
		return
	}

	hasDeleted, err := o.DeleteAllPortfolios()
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if !hasDeleted {
		writeError(w, http.StatusNotFound, fmt.Sprint("No portfolios to delete"))
		return
	}

	writeOk(w, struct {
		Status     responseStatus `json:"status"`
		HasDeleted bool           `json:"hasDeleted"`
	}{Status: ok, HasDeleted: true})
}
