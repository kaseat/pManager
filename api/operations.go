package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kaseat/pManager/portfolio"
)

// CreateSingleOperation creates single operation
// @summary Add new operation
// @description Adds operation to specified portfolio
// @id operation-add
// @accept json
// @produce json
// @param id path string true "Portfolio Id"
// @param portfolio body operationRequest true "Operation info"
// @success 200 {object} addOperationSuccess "Returns portfolio Id just created"
// @failure 400 {object} errorResponse "Returns when any processing error occurs"
// @failure 401 {object} errorResponse "Returns when authentication error occurs"
// @tags operations,portfolios
// @security ApiKeyAuth
// @router /portfolios/{id}/operations [post]
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

	writeOk(w, addOperationSuccess{OperationID: opID})
}

// ReadOperations gets all operations of specified portfolio
// @summary Get all operations
// @description Gets all operations for specified portfolio
// @id operation-get-all
// @produce json
// @param id path string true "Portfolio Id"
// @param figi query string false "Filter by FIGI"
// @param from query string false "Filter operations from this date"
// @param to query string false "Filter operations till this date"
// @success 200 {array} portfolio.Operation "Returns operations info"
// @failure 400 {object} errorResponse "Returns when any processing error occurs"
// @failure 401 {object} errorResponse "Returns when authentication error occurs"
// @tags operations,portfolios
// @security ApiKeyAuth
// @router /portfolios/{id}/operations [get]
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

	writeOk(w, ops)
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

	writeOk(w, struct {
		Status     responseStatus `json:"status"`
		NumDeleted int64          `json:"numDeleted"`
	}{Status: ok, NumDeleted: numDeleted})
}
