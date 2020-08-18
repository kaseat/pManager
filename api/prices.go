package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/kaseat/pManager/models"
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

// GetPrices gets prices
// @summary Get prices
// @description Get prices
// @id get-price
// @produce json
// @param isin query string false "ISIN"
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

// AddPrices gets prices
// @summary Add prices
// @description Get prices
// @id add-price
// @produce json
// @param isin query string false "ISIN"
// @param price body []priceRequest true "Price info"
// @success 200 {array} commonResponse "Returns success status"
// @failure 400 {object} errorResponse "Returns when any processing error occurs"
// @failure 401 {object} errorResponse "Returns when authentication error occurs"
// @tags prices
// @security ApiKeyAuth
// @router /prices [post]
func AddPrices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	var pricesRaw []priceRequest
	err = json.Unmarshal(body, &pricesRaw)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	s := storage.GetStorage()
	isin := r.FormValue("isin")
	prices := make([]models.Price, len(pricesRaw))
	for i, price := range pricesRaw {
		prices[i] = models.Price{
			Price:  price.Price,
			Volume: price.Volume,
			ISIN:   isin,
			Date:   time.Unix(price.Time, 0),
		}
	}

	err = s.AddPrices(prices)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeOk(w, commonResponse{Status: "ok"})
}
