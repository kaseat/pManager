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
// @summary Add new portfolio
// @description Creates single portfolio
// @id portfolio-add
// @accept json
// @produce json
// @param portfolio body portfolioRequest true "Portfolio info"
// @success 200 {object} addPortfoliioSuccess "Returns portfolio Id just created"
// @failure 400 {object} errorResponse "Returns when any processing error occurs"
// @failure 401 {object} errorResponse "Returns when authentication error occurs"
// @tags portfolios
// @security ApiKeyAuth
// @router /portfolios [post]
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

	writeOk(w, addPortfoliioSuccess{PortfolioID: p.PortfolioID})
}

// ReadSinglePortfolio gets single portfolio by id
// @summary Get portfolio by Id
// @description Gets portfolio info by Id
// @id portfolio-get-by-id
// @produce json
// @param id path string true "Portfolio Id"
// @success 200 {object} portfolio.Portfolio "Returns portfolio info if any"
// @failure 400 {object} errorResponse "Returns when any processing error occurs"
// @failure 401 {object} errorResponse "Returns when authentication error occurs"
// @tags portfolios
// @security ApiKeyAuth
// @router /portfolios/{id} [get]
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

	writeOk(w, p)
}

// ReadAllPortfolios gets all portfolios
// @summary Get all portfolios
// @description Gets all portfolios avaliable
// @id portfolio-get-all
// @produce json
// @success 200 {array} portfolio.Portfolio "Returns portfolio info"
// @failure 400 {object} errorResponse "Returns when any processing error occurs"
// @failure 401 {object} errorResponse "Returns when authentication error occurs"
// @tags portfolios
// @security ApiKeyAuth
// @router /portfolios [get]
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
		ps = []portfolio.Portfolio{}
	}

	writeOk(w, ps)
}

// UptateSinglePortfolio updates single portfolio by id
// @summary Update portfolio info
// @description Updates portfolio info by Id
// @id portfolio-put-by-id
// @accept json
// @produce json
// @param id path string true "Portfolio Id"
// @param portfolio body portfolioRequest true "Portfolio info"
// @success 200 {object} putPortfoliioSuccess "Returns portfolio info if any"
// @failure 400 {object} errorResponse "Returns when any processing error occurs"
// @failure 401 {object} errorResponse "Returns when authentication error occurs"
// @tags portfolios
// @security ApiKeyAuth
// @router /portfolios/{id} [put]
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

	p.PortfolioID = mux.Vars(r)["id"]
	p.OwnerID = o.OwnerID

	hasUptated, err := p.UpdatePortfolio()
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeOk(w, putPortfoliioSuccess{HasModified: hasUptated})
}

// DeleteSinglePortfolio deletes single portfolio by id
// @summary Delete portfolio
// @description Deletes portfolio an all associated operations
// @id portfolio-del-by-id
// @produce json
// @param id path string true "Portfolio Id"
// @success 200 {object} delPortfoliioSuccess "Returns true if portfolio has deleted"
// @failure 400 {object} errorResponse "Returns when any processing error occurs"
// @failure 401 {object} errorResponse "Returns when authentication error occurs"
// @tags portfolios
// @security ApiKeyAuth
// @router /portfolios/{id} [delete]
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

	writeOk(w, delPortfoliioSuccess{HasDeleted: hasDeleted})
}

// DeleteAllPortfolios deletes all portfolios
// @summary Delete all portfolios
// @description Deletes all portfolios an all associated operations
// @id portfolio-del-all
// @produce json
// @success 200 {object} delPortfoliioSuccess "Returns true if portfolios has deleted"
// @failure 400 {object} errorResponse "Returns when any processing error occurs"
// @failure 401 {object} errorResponse "Returns when authentication error occurs"
// @tags portfolios
// @security ApiKeyAuth
// @router /portfolios [delete]
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

	writeOk(w, delPortfoliioSuccess{HasDeleted: hasDeleted})
}
