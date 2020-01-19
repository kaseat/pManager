package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kaseat/tcssync/portfolio"
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

type portfolioCreatedResponse struct {
	Status      responseStatus `json:"status"`
	PortfolioID string         `json:"createdPortfolioId"`
}
type portfolioGetAllResponse struct {
	Status     responseStatus        `json:"status"`
	Portfolios []portfolio.Portfolio `json:"portfolios"`
}

func main() {
	cfg := portfolio.Config{
		MongoURL: "mongodb://localhost:27017",
		DbName:   "tcs",
	}
	portfolio.Init(cfg)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/portfolio", addPortfolio).Methods("POST")
	router.HandleFunc("/portfolio", getAllPortfolios).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getAllPortfolios(w http.ResponseWriter, r *http.Request) {
	ps, err := portfolio.GetAllPortfolios()
	if err != nil {
		resp := errorResponse{Status: notOk, Error: err.Error()}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&resp)
		return
	}

	resp := portfolioGetAllResponse{Status: ok, Portfolios: ps}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&resp)
}

func addPortfolio(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp := errorResponse{Status: notOk, Error: err.Error()}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&resp)
		return
	}

	var p portfolio.Portfolio
	err = json.Unmarshal(body, &p)
	if err != nil {
		resp := errorResponse{Status: notOk, Error: err.Error()}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&resp)
		return
	}

	if p.Name == "" {
		errText := "You must provide portfolio name"
		resp := errorResponse{Status: notOk, Error: errText}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&resp)
		return
	}

	p, err = portfolio.AddPortfolio(p.Name, p.Description)
	if err != nil {
		resp := errorResponse{Status: notOk, Error: err.Error()}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&resp)
		return
	}

	w.WriteHeader(http.StatusCreated)
	resp := portfolioCreatedResponse{Status: ok, PortfolioID: p.PortfolioID}
	json.NewEncoder(w).Encode(&resp)
}
