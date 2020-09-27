package moex

import (
	"errors"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/kaseat/pManager/models"
	"github.com/kaseat/pManager/storage"
)

var isSync int32

// Sync starts moex sync
func Sync(ticker string) {
	defer func() {
		atomic.StoreInt32(&isSync, 0)
		fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "End sync MOEX")
	}()
	fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Begin sync MOEX")
	if atomic.LoadInt32(&isSync) == 1 {
		err := errors.New("MOEX sync already in process")
		fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Error sync instruments:", err)
		return
	}
	atomic.StoreInt32(&isSync, 1)

	client := &http.Client{}

	s := storage.GetStorage()

	var instrumentsToSync []models.Instrument
	var err error = nil
	if ticker == "" {
		instrumentsToSync, err = s.GetInstruments("code", "MOEX")
	} else {
		instrumentsToSync, err = s.GetInstruments("ticker", ticker)
	}

	if err != nil {
		fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Error get securities list:", err)
		return
	}

	for _, inst := range instrumentsToSync {
		from := inst.PriceUptdTime
		if from.IsZero() {
			from = time.Date(2019, time.May, 1, 0, 0, 0, 0, time.UTC)
		}
		var pricesRaw []priceInternal

		for cur := 0; ; {
			pr, curRe, err := fetchFromAPI(client, from, inst.Ticker, cur)
			if err != nil {
				break
			}
			pricesRaw = append(pricesRaw, pr...)
			if curRe == 0 {
				break
			}
			cur = curRe
		}
		prices := make([]models.Price, 0, len(pricesRaw))
		var lastDate time.Time
		for _, priceRaw := range pricesRaw {
			if inst.Currency != priceRaw.Currency {
				continue
			}
			price := models.Price{
				SecID:  inst.SecID,
				Price:  priceRaw.Price,
				Volume: priceRaw.Volume,
				Date:   priceRaw.Date,
				ISIN:   inst.ISIN,
			}
			if priceRaw.Date.After(lastDate) {
				lastDate = priceRaw.Date
			}
			prices = append(prices, price)
		}
		if len(prices) > 0 {
			if err = s.AddPrices(prices); err != nil {
				fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Error add", len(prices), "prices for", ticker, "to storage:", err)
			} else {
				fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Success add", len(prices), "prices for", ticker, "to storage")
				s.SetInstrumentPriceUptdTime(inst.SecID, lastDate.AddDate(0, 0, 1))
			}
		}
	}
}
