package api

import (
	"net/http"

	"github.com/kaseat/pManager/storage"
	"github.com/kaseat/pManager/sync/tcs"
)

// SyncPrices sync prices
// @summary Sync prices
// @description Sync prices
// @id sync-price
// @produce json
// @success 200 {array} commonResponse "Returns success status"
// @failure 400 {object} errorResponse "Returns when any processing error occurs"
// @failure 401 {object} errorResponse "Returns when authentication error occurs"
// @tags prices
// @security ApiKeyAuth
// @router /prices/sync [get]
func SyncPrices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	go tcs.SyncPrices()
	writeOk(w, commonResponse{Status: "ok"})
}

// GetPrices sync prices
// @summary Get prices
// @description Sync prices
// @id sync-price
// @produce json
// @param from query string false "Filter prices from this date"
// @param to query string false "Filter prices till this date"
// @success 200 {array} commonResponse "Returns success status"
// @failure 400 {object} errorResponse "Returns when any processing error occurs"
// @failure 401 {object} errorResponse "Returns when authentication error occurs"
// @tags prices
// @security ApiKeyAuth
// @router /prices [get]
func GetPrices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	s := storage.GetStorage()
	isin := r.FormValue("isin")
	from := r.FormValue("from")
	to := r.FormValue("to")

	prices, err := s.GetPricesByIsin(isin, from, to)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeOk(w, prices)
}
