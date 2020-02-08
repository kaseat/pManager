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

	p, err = portfolio.AddPortfolio(p.Name, p.Description)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp := singlePortfolioIDResponse{Status: ok, PortfolioID: p.PortfolioID}
	bytes, _ := json.Marshal(&resp)

	w.WriteHeader(http.StatusCreated)
	w.Write(bytes)
}

// ReadSinglePortfolio gets single portfolio by id
func ReadSinglePortfolio(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	found, p, err := portfolio.GetPortfolio(id)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if !found {
		writeError(w, http.StatusNotFound, fmt.Sprint("No portfolio found with Id: ", id))
		return
	}

	resp := singlePortfolioResponse{Status: ok, Portfolio: p}
	bytes, _ := json.Marshal(&resp)

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

// ReadAllPortfolios gets all portfolios
func ReadAllPortfolios(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ps, err := portfolio.GetAllPortfolios()
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(ps) == 0 {
		writeError(w, http.StatusNotFound, fmt.Sprint("No portfolios found"))
		return
	}

	resp := multiplePortfoliosResponse{Status: ok, Portfolios: ps}
	bytes, _ := json.Marshal(&resp)

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
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

	resp := singlePortfolioUpdateResponse{Status: ok, HasUptated: hasUptated}
	bytes, _ := json.Marshal(&resp)

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

// DeleteSinglePortfolio deletes single portfolio by id
func DeleteSinglePortfolio(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	hasDeleted, err := portfolio.DeletePortfolio(id)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if !hasDeleted {
		writeError(w, http.StatusNotFound, fmt.Sprint("No portfolio found with Id: ", id))
		return
	}

	resp := portfolioDeleteResponse{Status: ok, HasDeleted: true}
	bytes, _ := json.Marshal(&resp)

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

// DeleteAllPortfolios deletes all portfolios
func DeleteAllPortfolios(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	hasDeleted, err := portfolio.DeleteAllPortfolios()
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if !hasDeleted {
		writeError(w, http.StatusNotFound, fmt.Sprint("No portfolios to delete"))
		return
	}

	resp := portfolioDeleteResponse{Status: ok, HasDeleted: true}
	bytes, _ := json.Marshal(&resp)

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
