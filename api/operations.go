package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kaseat/pManager/portfolio"
)

// CreateSingleOperation creates single operation
func CreateSingleOperation(w http.ResponseWriter, r *http.Request) {
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

	decoder := json.NewDecoder(r.Body)
	var op portfolio.Operation

	err = decoder.Decode(&op)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	opID, err := p.AddOperation(op)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp := singleOperationIDResponse{Status: ok, OperationID: opID}
	bytes, _ := json.Marshal(&resp)

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

// ReadAllOperations gets all operations of specified portfolio
func ReadAllOperations(w http.ResponseWriter, r *http.Request) {
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

	figi := r.FormValue("figi")
	var ops []portfolio.Operation
	if figi != "" {
		ops, err = p.GetAllOperationsByFigi(figi)
	} else {
		ops, err = p.GetAllOperations()
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if ops == nil {
		ops = []portfolio.Operation{}
	}

	for i := range ops {
		ops[i].PortfolioID = ""
	}

	resp := operationsResponse{Status: ok, PortfolioID: p.PortfolioID, Operations: ops}
	bytes, _ := json.Marshal(&resp)

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

// DeleteAllOperations removes all operations of specified portfolio
func DeleteAllOperations(w http.ResponseWriter, r *http.Request) {
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

	numDeleted, err := p.DeleteAllOperations()
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp := operationDeleteResponse{Status: ok, NumDeleted: numDeleted}
	bytes, _ := json.Marshal(&resp)

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
