package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kaseat/pManager/models"
	"github.com/kaseat/pManager/storage"
	"github.com/kaseat/pManager/utils"
)

// GetAveragePrice returns averge price of specified ticker
// @summary Get average
// @description Gets average price of given ticker
// @id get-average
// @produce json
// @param id path string true "Portfolio Id"
// @param ticker query string true "Ticker"
// @success 200 {array} getAverageSuccess "Returns average price of given ticker"
// @failure 400 {object} errorResponse "Returns when any processing error occurs"
// @failure 401 {object} errorResponse "Returns when authentication error occurs"
// @tags misc
// @security ApiKeyAuth
// @router /portfolios/{id}/average [get]
func GetAveragePrice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pid := mux.Vars(r)["id"]
	user := r.Header.Get("user")

	s := storage.GetStorage()

	canAccess, err := canAccess(s, user, pid)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if !canAccess {
		writeError(w, http.StatusUnauthorized, "You cannot get operations from this portfolio")
		return
	}

	ticker := r.FormValue("ticker")
	if ticker == "" {
		writeError(w, http.StatusBadRequest, fmt.Sprint(`You must provide "ticker" parameter`))
		return
	}

	ops, err := s.GetOperations(pid, "ticker", ticker, "", "")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	avg := utils.GetAverage(ops)

	writeOk(w, getAverageSuccess{Average: avg})
}

// GetBalance returns balance of specified currency
// @summary Get balance
// @description Gets balance of given currency
// @id get-balance
// @produce json
// @param id path string true "Portfolio Id"
// @param currency query string true "Currency"
// @success 200 {array} getBalanceSuccess "Returns balance of given currency"
// @failure 400 {object} errorResponse "Returns when any processing error occurs"
// @failure 401 {object} errorResponse "Returns when authentication error occurs"
// @tags misc
// @security ApiKeyAuth
// @router /portfolios/{id}/balance [get]
func GetBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pid := mux.Vars(r)["id"]
	user := r.Header.Get("user")

	curr := models.Currency(r.FormValue("currency"))
	if curr == "" {
		writeError(w, http.StatusBadRequest, fmt.Sprint("You must provide 'currency' parameter"))
		return
	}

	if !(curr == models.EUR || curr == models.RUB || curr == models.USD) {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("Unknown currency '%s'. Expected '%s', '%s' or '%s'",
			curr, models.EUR, models.RUB, models.USD))
		return
	}

	s := storage.GetStorage()

	canAccess, err := canAccess(s, user, pid)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if !canAccess {
		writeError(w, http.StatusUnauthorized, "You cannot get operations from this portfolio")
		return
	}

	ops, err := s.GetOperations(pid, "curr", string(curr), "", "")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	bal := utils.GetSum(ops)

	writeOk(w, getBalanceSuccess{Balance: bal})
}
