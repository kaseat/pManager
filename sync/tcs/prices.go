package tcs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/kaseat/pManager/models"
	"github.com/kaseat/pManager/storage"
)

const candlesURL = "https://api-invest.tinkoff.ru/openapi/sandbox/market/candles"

// SyncPrices sync daily prices for instruments
func SyncPrices() {
	s := storage.GetStorage()
	token := s.GetTcsToken()
	instruments, _ := s.GetAllInstruments()
	client := &http.Client{}

	for _, x := range instruments {
		beginDate := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := today()
		chunks := getTimeChunks(beginDate, endDate, 12)
		for _, ch := range chunks {
			time.Sleep(500 * time.Millisecond)
			if err := s.AddPrices(getPrices(client, token, x, ch.From, ch.To)); err != nil {
				fmt.Printf("Sync price for %s from %s to %s error: %v\n", x.Ticker, ch.From.Format("2006-01-02"), ch.To.Format("2006-01-02"), err)
			} else {
				fmt.Printf("Sync price for %s from %s to %s succeded\n", x.Ticker, ch.From.Format("2006-01-02"), ch.To.Format("2006-01-02"))
			}
		}
	}

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
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(body, &respObj)
	if err != nil {
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

func today() time.Time {
	year, month, day := time.Now().Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func getTimeChunks(from, to time.Time, size int) []chunk {
	result := []chunk{}
	for {
		end := from.AddDate(0, size, 0)
		if to.Before(end) {
			result = append(result, chunk{From: from, To: to})
			break
		}
		result = append(result, chunk{From: from, To: end})
		from = end.AddDate(0, 0, 1)
	}
	return result
}

type chunk struct {
	From time.Time
	To   time.Time
}
