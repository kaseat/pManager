package tcs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/kaseat/pManager/models"
	"github.com/kaseat/pManager/storage"
)

const candlesURL = "https://api-invest.tinkoff.ru/openapi/sandbox/market/candles"

var lastSyncPricesError atomic.Value
var syncPricesIsRunning int32

// SyncPrices sync daily prices for prices
func SyncPrices() {
	defer atomic.StoreInt32(&syncPricesIsRunning, 0)
	if atomic.LoadInt32(&syncPricesIsRunning) == 1 {
		return
	}
	atomic.StoreInt32(&syncPricesIsRunning, 1)
	fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Begin sync prices")
	s := storage.GetStorage()
	token, _ := s.GetTcsToken()
	instruments, _ := s.GetAllInstruments()
	client := &http.Client{}

	for _, x := range instruments {
		beginDate := x.PriceUptdTime
		if beginDate.IsZero() {
			beginDate = time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
		}

		endDate := today()
		chunks := getTimeChunks(beginDate, endDate, 12)
		for _, ch := range chunks {
			if ch.From == ch.To {
				break
			}
			time.Sleep(500 * time.Millisecond)
			now := time.Now().Format("2006-02-01 15:04:05")
			from := ch.From.Format("2006-01-02")
			to := ch.To.Format("2006-01-02")
			if err := s.AddPrices(getPrices(client, token, x, ch.From, ch.To)); err != nil {
				fmt.Printf("%s Sync price for %s from %s to %s error: %v\n", now, x.Ticker, from, to, err)
			} else {
				s.SetInstrumentPriceUptdTime(x.ISIN, ch.To)
				fmt.Printf("%s Sync price for %s from %s to %s succeded\n", now, x.Ticker, from, to)
			}
		}
	}
	fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Success sync prices")
}

func getPrices(client *http.Client, token string, ins models.Instrument, from, to time.Time) []models.Price {
	var respObj struct {
		Payload struct {
			Candles []struct {
				Price float64   `json:"c"`
				Vol   int       `json:"v"`
				Time  time.Time `json:"time"`
			} `json:"candles"`
		} `json:"payload"`
	}

	req, err := http.NewRequest("GET", candlesURL, nil)
	if err != nil {
		setLastPricesError(err)
		return nil
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.URL.RawQuery = url.Values{
		"figi":     {ins.FIGI},
		"from":     {from.Format("2006-01-02T15:04:05Z")},
		"to":       {to.Format("2006-01-02T15:04:05Z")},
		"interval": {"day"},
	}.Encode()

	resp, err := client.Do(req)
	if err != nil {
		setLastPricesError(err)
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		setLastPricesError(err)
		return nil
	}
	err = json.Unmarshal(body, &respObj)
	if err != nil {
		setLastPricesError(err)
		return nil
	}

	result := make([]models.Price, len(respObj.Payload.Candles))
	for i, x := range respObj.Payload.Candles {
		result[i] = models.Price{
			Date:   x.Time,
			ISIN:   ins.ISIN,
			Price:  x.Price,
			Volume: x.Vol,
		}
	}

	return result
}

func setLastPricesError(err error) {
	fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Error sync instruments:", err)
	lastSyncPricesError.Store(syncError{Error: err, IsNotEmpty: true})
}
