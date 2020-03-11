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

// ReadOperations gets all operations of specified portfolio
func ReadOperations(w http.ResponseWriter, r *http.Request) {
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

	ops, err := p.GetOperations(r.FormValue("figi"), r.FormValue("from"), r.FormValue("to"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	for i := range ops {
		ops[i].PortfolioID = ""
		ops[i].PriceF = float64(ops[i].Price) / 1e6
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
