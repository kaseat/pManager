package api

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kaseat/pManager/portfolio"
)

// GetAveragePrice returns averge price of specified figi
func GetAveragePrice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	found, p, err := portfolio.GetPortfolio(id)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if !found {
		writeError(w, http.StatusNotFound, fmt.Sprint("No portfolio found with Id: ", id))
		return
	}
	figi := mux.Vars(r)["figi"]

	var price float64
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
	price = math.Round(price*100) / 100
	resp := struct {
		Status responseStatus
		Price  float64
	}{
		Status: ok,
		Price:  price,
	}

	bytes, _ := json.Marshal(&resp)

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
