package main

import (
	"encoding/json"
	"fmt"
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

type singlePortfolioIDResponse struct {
	Status      responseStatus `json:"status"`
	PortfolioID string         `json:"createdPortfolioId"`
}

type singlePortfolioUpdateResponse struct {
	Status     responseStatus `json:"status"`
	HasUptated bool           `json:"hasModified"`
}

type portfolioDeleteResponse struct {
	Status     responseStatus `json:"status"`
	HasDeleted bool           `json:"hasDeleted"`
}

type multiplePortfoliosResponse struct {
	Status     responseStatus        `json:"status"`
	Portfolios []portfolio.Portfolio `json:"portfolios"`
}

type singlePortfolioResponse struct {
	Status    responseStatus      `json:"status"`
	Portfolio portfolio.Portfolio `json:"portfolio"`
}

type operationsResponse struct {
	Status     responseStatus        `json:"status"`
	Portfolio  portfolio.Portfolio   `json:"portfolio"`
	Operations []portfolio.Operation `json:"operations"`
}

func main() {
	cfg := portfolio.Config{
		MongoURL: "mongodb://localhost:27017",
		DbName:   "tcs",
	}
	portfolio.Init(cfg)

	fmt.Println("Started!")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/portfolios", createSinglePortfolio).Methods("POST")
	router.HandleFunc("/portfolios", readAllPortfolios).Methods("GET")
	router.HandleFunc("/portfolios", deleteAllPortfolios).Methods("DELETE")
	router.HandleFunc("/portfolios/{id}", readSinglePortfolio).Methods("GET")
	router.HandleFunc("/portfolios/{id}", uptateSinglePortfolio).Methods("PUT")
	router.HandleFunc("/portfolios/{id}", deleteSinglePortfolio).Methods("DELETE")
	router.HandleFunc("/portfolios/{id}/operations", readAllOperations).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func readAllOperations(w http.ResponseWriter, r *http.Request) {
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

	ops, err := p.GetAllOperations()
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if ops == nil {
		ops = []portfolio.Operation{}
	}

	resp := operationsResponse{Status: ok, Portfolio: p, Operations: ops}
	bytes, _ := json.Marshal(&resp)

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func createSinglePortfolio(w http.ResponseWriter, r *http.Request) {
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

func readSinglePortfolio(w http.ResponseWriter, r *http.Request) {
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

func readAllPortfolios(w http.ResponseWriter, r *http.Request) {
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

func uptateSinglePortfolio(w http.ResponseWriter, r *http.Request) {
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

func deleteSinglePortfolio(w http.ResponseWriter, r *http.Request) {
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

func deleteAllPortfolios(w http.ResponseWriter, r *http.Request) {
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

func writeError(w http.ResponseWriter, statusCode int, text string) {
	resp := errorResponse{Status: notOk, Error: text}
	bytes, _ := json.Marshal(&resp)
	w.WriteHeader(statusCode)
	w.Write(bytes)
}
