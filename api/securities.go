package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kaseat/pManager/models"
	"github.com/kaseat/pManager/storage"
	"github.com/kaseat/pManager/sync/tcs"
)

// SyncSecurities syncs securities
// @summary Sync securities
// @description Sync intruments dimension
// @id sync-securities
// @produce json
// @success 200 {array} commonResponse "Returns success status"
// @failure 400 {object} errorResponse "Returns when any processing error occurs"
// @failure 401 {object} errorResponse "Returns when authentication error occurs"
// @tags securities
// @security ApiKeyAuth
// @router /securities/sync [get]
func SyncSecurities(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stat := tcs.GetSyncInstrumentsStatus()
	if stat.Status != tcs.Processing {
		go tcs.SyncInstruments()
		writeOk(w, commonResponse{Status: "ok"})
	} else {
		writeError(w, http.StatusBadRequest, "Sync already in process")
	}
}

// GetSecurities gets securities
// @summary Get securities
// @description Gets securities avaliable
// @id get-securities
// @produce json
// @param filter query string false "Filter by" Enums(none, ticker, isin, figi) default(none)
// @param by query string false "Filter value"
// @success 200 {array} models.Instrument "Returns success status"
// @failure 400 {object} errorResponse "Returns when any processing error occurs"
// @failure 401 {object} errorResponse "Returns when authentication error occurs"
// @tags securities
// @security ApiKeyAuth
// @router /securities [get]
func GetSecurities(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	s := storage.GetStorage()
	filter := r.FormValue("filter")
	by := r.FormValue("by")

	var err error
	var ins []models.Instrument
	if filter == "none" || filter == "" {
		ins, err = s.GetAllInstruments()
	} else {
		ins, err = s.GetInstruments(filter, by)
	}

	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeOk(w, ins)
}

// AddSecurities adds securities
// @summary Add securities
// @description Adds securities
// @id add-securities
// @produce json
// @param instrument body models.Instrument true "Instrument info"
// @success 200 {array} models.Instrument "Returns success status"
// @failure 400 {object} errorResponse "Returns when any processing error occurs"
// @failure 401 {object} errorResponse "Returns when authentication error occurs"
// @tags securities
// @security ApiKeyAuth
// @router /securities [post]
func AddSecurities(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var ins models.Instrument

	err := decoder.Decode(&ins)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	s := storage.GetStorage()
	err = s.AddInstruments([]models.Instrument{ins})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
	} else {
		writeOk(w, commonResponse{Status: "ok"})
	}
}

// GetSecuritiesForPortfolio gets securities for given portfolio
// @summary Get securities for given portfolio
// @description Gets portfolio info by Id
// @id get-by-portfolio-securities
// @produce json
// @param id path string true "Portfolio Id"
// @param on query string false "Get securities on this date"
// @success 200 {object} models.Portfolio "Returns portfolio info if any"
// @failure 400 {object} errorResponse "Returns when any processing error occurs"
// @failure 401 {object} errorResponse "Returns when authentication error occurs"
// @tags securities
// @security ApiKeyAuth
// @router /portfolios/{id}/securities [get]
func GetSecuritiesForPortfolio(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pid := mux.Vars(r)["id"]
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

	canRead := false
	for _, p := range ps {
		if p.PortfolioID == pid {
			canRead = true
			break
		}
	}

	if !canRead {
		writeError(w, http.StatusUnauthorized, "You cannot read securities for this portfolio")
		return
	}

	ops, err := s.GetShares(pid, r.FormValue("on"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeOk(w, ops)
}
