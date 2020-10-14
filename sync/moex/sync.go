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
func Sync(ticker string, httpClient *http.Client) {
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

	today := today()
	for _, instrument := range instrumentsToSync {
		securityInfo, err := getSecurityInfo(httpClient, instrument.Ticker)
		if err != nil {
			fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Error get instrument info:", err)
			continue
		}

		from := instrument.PriceUptdTime
		if from.IsZero() {
			from = time.Date(2019, time.May, 1, 0, 0, 0, 0, time.UTC)
		}

		if from.After(today) {
			continue
		}

		var pricesRaw []priceInternal

		for cursor := 0; ; {
			rawPrice, cursorAfterFetch, err := fetchFromAPI(httpClient, from, securityInfo, cursor)
			if err != nil {
				break
			}
			pricesRaw = append(pricesRaw, rawPrice...)
			if cursorAfterFetch == 0 {
				break
			}
			cursor = cursorAfterFetch
		}
		prices := make([]models.Price, 0, len(pricesRaw))
		var lastDate time.Time
		for _, rawPrice := range pricesRaw {
			if instrument.Currency != rawPrice.Currency {
				continue
			}
			price := models.Price{
				SecID:  instrument.SecID,
				Price:  rawPrice.Price,
				Volume: rawPrice.Volume,
				Date:   rawPrice.Date,
				ISIN:   instrument.ISIN,
			}
			if rawPrice.Date.After(lastDate) {
				lastDate = rawPrice.Date
			}
			prices = append(prices, price)
		}
		if len(prices) > 0 {
			if err = s.AddPrices(prices); err != nil {
				fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Error add", len(prices), "prices for", ticker, "to storage:", err)
			} else {
				fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Success add", len(prices), "prices for", ticker, "to storage")
				s.SetInstrumentPriceUptdTime(instrument.SecID, lastDate.AddDate(0, 0, 1))
			}
		}
	}
}

func today() time.Time {
	now := time.Now().UTC()
	y, m, d := now.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, now.Location()).Add(time.Hour * -1)
}
