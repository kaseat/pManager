package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kaseat/pManager/models"
	"github.com/kaseat/pManager/storage"
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

	s := storage.GetStorage()
	u, err := s.GetUserByLogin(user)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var p models.Portfolio
	err = json.Unmarshal(body, &p)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if p.Name == "" {
		writeError(w, http.StatusBadRequest, "You must provide portfolio name")
		return
	}

	pid, err := s.AddPortfolio(u.UserID, p)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeOk(w, addPortfoliioSuccess{PortfolioID: pid})
}

// ReadSinglePortfolio gets single portfolio by id
// @summary Get portfolio by Id
// @description Gets portfolio info by Id
// @id portfolio-get-by-id
// @produce json
// @param id path string true "Portfolio Id"
// @success 200 {object} models.Portfolio "Returns portfolio info if any"
// @failure 400 {object} errorResponse "Returns when any processing error occurs"
// @failure 401 {object} errorResponse "Returns when authentication error occurs"
// @tags portfolios
// @security ApiKeyAuth
// @router /portfolios/{id} [get]
func ReadSinglePortfolio(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	user := r.Header.Get("user")

	s := storage.GetStorage()
	u, err := s.GetUserByLogin(user)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	p, err := s.GetPortfolio(u.UserID, id)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeOk(w, p)
}

// ReadAllPortfolios gets all portfolios
// @summary Get all portfolios
// @description Gets all portfolios avaliable
// @id portfolio-get-all
// @produce json
// @success 200 {array} models.Portfolio "Returns portfolio info"
// @failure 400 {object} errorResponse "Returns when any processing error occurs"
// @failure 401 {object} errorResponse "Returns when authentication error occurs"
// @tags portfolios
// @security ApiKeyAuth
// @router /portfolios [get]
func ReadAllPortfolios(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user := r.Header.Get("user")

	s := storage.GetStorage()
	u, err := s.GetUserByLogin(user)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	ps, err := s.GetPortfolios(u.UserID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
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
	var p models.Portfolio

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
	pid := mux.Vars(r)["id"]

	s := storage.GetStorage()
	u, err := s.GetUserByLogin(user)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	modified, err := s.UpdatePortfolio(u.UserID, pid, p)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeOk(w, putPortfoliioSuccess{HasModified: modified})
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

	s := storage.GetStorage()
	u, err := s.GetUserByLogin(user)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	deleted, err := s.DeletePortfolio(u.UserID, id)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeOk(w, delPortfoliioSuccess{HasDeleted: deleted})
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

	s := storage.GetStorage()
	u, err := s.GetUserByLogin(user)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	num, err := s.DeletePortfolios(u.UserID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeOk(w, delPortfoliioSuccess{HasDeleted: num > 0})
}
