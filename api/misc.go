package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kaseat/pManager/models"
	"github.com/kaseat/pManager/portfolio"
	"github.com/kaseat/pManager/storage"
	"github.com/kaseat/pManager/utils"
)

// GetAveragePrice returns averge price of specified figi
func GetAveragePrice(w http.ResponseWriter, r *http.Request) {
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
	figi := r.FormValue("figi")
	if figi == "" {
		writeError(w, http.StatusBadRequest, fmt.Sprint(`You must provide "figi" parameter`))
		return
	}

	var price int64
	onDate, err := time.Parse("2006-01-02T15:04:05.000Z0700", r.FormValue("onDate"))
	if err == nil {
		price, err = p.GetAveragePriceByFigiTillDate(figi, onDate)
	} else {
		price, err = p.GetAveragePriceByFigi(figi)
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeOk(w, struct {
		Status responseStatus
		Price  float64
	}{
		Status: ok,
		Price:  float64(price) / 1e6,
	})
}

// GetBalance returns balance of specified currency
// @summary Get all operations
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
			curr, portfolio.EUR, portfolio.RUB, portfolio.USD))
		return
	}

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

	canDel := false
	for _, p := range ps {
		if p.PortfolioID == pid {
			canDel = true
			break
		}
	}

	if !canDel {
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
