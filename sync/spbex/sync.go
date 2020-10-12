package spbex

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/kaseat/pManager/models"
	"github.com/kaseat/pManager/storage"
)

var isSync int32

// Sync starts spbex sync
func Sync(ticker string, httpClient *http.Client) {
	defer func() {
		atomic.StoreInt32(&isSync, 0)
		fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "End sync SPBEX")
	}()
	fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Begin sync SPBEX")
	if atomic.LoadInt32(&isSync) == 1 {
		err := errors.New("SPBEX sync already in process")
		fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Error sync prices:", err)
		return
	}
	atomic.StoreInt32(&isSync, 1)

	s := storage.GetStorage()

	var instrumentsToSync []models.Instrument
	var err error = nil
	if ticker == "" {
		instrumentsToSync, err = s.GetInstruments("code", "SPBEX")
	} else {
		instrumentsToSync, err = s.GetInstruments("ticker", ticker)
	}

	if err != nil {
		fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Error get securities list:", err)
		return
	}

	for _, instrument := range instrumentsToSync {

		from := instrument.PriceUptdTime
		if from.IsZero() {
			from = time.Date(2019, time.May, 1, 0, 0, 0, 0, time.UTC)
		}

		url := "https://investcab.ru/api/chistory?symbol=%s&resolution=D"
		url = fmt.Sprintf(url, instrument.Ticker)
		url += fmt.Sprintf("&from=%d&to=%d", from.Unix(), time.Now().Unix())

		resp, err := httpClient.Get(url)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()

		r, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		var rawPrices struct {
			Timestamp []int     `json:"t"`
			Price     []float64 `json:"c"`
		}
		err = json.Unmarshal([]byte(strings.Replace(string(r)[1:len(r)-1], `\"`, `"`, -1)), &rawPrices)
		if err != nil {
			fmt.Println(err)
			return
		}

		pricesInt := make([]struct {
			Date  time.Time
			Price float64
		}, len(rawPrices.Price))

		for i := 0; i < len(rawPrices.Price); i++ {
			pricesInt[i] = struct {
				Date  time.Time
				Price float64
			}{
				Date:  time.Unix(int64(rawPrices.Timestamp[i]), 0),
				Price: rawPrices.Price[i],
			}
		}

		prices := make([]models.Price, 0, len(pricesInt))
		var lastDate time.Time
		for _, rawPrice := range pricesInt {
			price := models.Price{
				SecID:  instrument.SecID,
				Price:  rawPrice.Price,
				Volume: 0,
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
				fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Error add", len(prices), "prices for", instrument.Ticker, "to storage:", err)
			} else {
				fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Success add", len(prices), "prices for", instrument.Ticker, "to storage")
				s.SetInstrumentPriceUptdTime(instrument.SecID, lastDate.AddDate(0, 0, 1))
			}
		}
	}
}
