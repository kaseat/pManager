package api

import (
	"encoding/json"
	"net/http"

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
