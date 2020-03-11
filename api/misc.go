package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kaseat/pManager/portfolio"
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
	fmt.Println(onDate)
	resp := struct {
		Status responseStatus
		Price  float64
	}{
		Status: ok,
		Price:  float64(price) / 1e6,
	}

	bytes, err := json.Marshal(&resp)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

// GetBalance returns balance of specified currency
func GetBalance(w http.ResponseWriter, r *http.Request) {
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
	curr := portfolio.Currency(r.FormValue("currency"))
	if curr == "" {
		writeError(w, http.StatusBadRequest, fmt.Sprint("You must provide 'currency' parameter"))
		return
	}
	if !(curr == portfolio.EUR || curr == portfolio.RUB || curr == portfolio.USD) {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("Unknown currency '%s'. Expected '%s', '%s' or '%s'",
			curr, portfolio.EUR, portfolio.RUB, portfolio.USD))
		return
	}

	balance, err := p.GetBalance(curr, "", r.FormValue("on"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	resp := struct {
		Status  responseStatus
		Balance float64
	}{
		Status:  ok,
		Balance: float64(balance) / 1e6,
	}

	bytes, err := json.Marshal(&resp)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
