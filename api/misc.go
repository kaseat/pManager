package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kaseat/pManager/gmail"
	"github.com/kaseat/pManager/models/currency"
	"github.com/kaseat/pManager/storage"
	"github.com/kaseat/pManager/sync"
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

	curr := currency.Type(r.FormValue("currency"))
	if curr == "" {
		writeError(w, http.StatusBadRequest, fmt.Sprint("You must provide 'currency' parameter"))
		return
	}

	if !(curr == currency.EUR || curr == currency.RUB || curr == currency.USD) {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("Unknown currency '%s'. Expected '%s', '%s' or '%s'",
			curr, currency.EUR, currency.RUB, currency.USD))
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

// AddGoogleAuth adds gmail support for import operations
// @summary Get GMail auth url
// @description Gets url for GMail auth
// @id get-gmail-url
// @produce json
// @success 200 {array} gmailAuthUrlSuccess "Returns url for GMail auth"
// @failure 400 {object} errorResponse "Returns when any processing error occurs"
// @failure 401 {object} errorResponse "Returns when authentication error occurs"
// @tags user
// @security ApiKeyAuth
// @router /misc/gmail/url [get]
func AddGoogleAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user := r.Header.Get("user")

	cl := gmail.GetClient()
	url, err := cl.GetAuthURL(user)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeOk(w, gmailAuthUrlSuccess{URL: url})
}

// AppCallback saves respose from gmail
func AppCallback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")

	cl := gmail.GetClient()
	err := cl.HandleResponse(state, code)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeOk(w, struct {
		Status string `json:"status"`
	}{Status: "ok"})
}

// SyncOperations sync operations
// @summary Sync operations
// @description Sync operations for given portfolio
// @id sync-op
// @produce json
// @param id path string true "Portfolio Id"
// @param from query string false "Filter operations from this date"
// @param to query string false "Filter operations till this date"
// @success 200 {array} commonResponse "Returns success status"
// @failure 400 {object} errorResponse "Returns when any processing error occurs"
// @failure 401 {object} errorResponse "Returns when authentication error occurs"
// @tags misc
// @security ApiKeyAuth
// @router /portfolios/{id}/sync [get]
func SyncOperations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	pid := mux.Vars(r)["id"]
	login := r.Header.Get("user")

	go sync.Sberbank(login, pid, r.FormValue("from"), r.FormValue("to"))

	writeOk(w, commonResponse{Status: "ok"})
}
