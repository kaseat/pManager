package tcs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/kaseat/pManager/models"
	"github.com/kaseat/pManager/storage"
)

const stocksURL = "https://api-invest.tinkoff.ru/openapi/sandbox/market/stocks"
const bondsURL = "https://api-invest.tinkoff.ru/openapi/sandbox/market/bonds"
const etfURL = "https://api-invest.tinkoff.ru/openapi/sandbox/market/etfs"

var lastSyncIstrumentsError error
var syncInstrumentsIsRunning int32

// SyncInstruments start sync instruments from tcs API
func SyncInstruments() {
	defer atomic.StoreInt32(&syncInstrumentsIsRunning, 0)
	if atomic.LoadInt32(&syncInstrumentsIsRunning) == 1 {
		return
	}
	atomic.StoreInt32(&syncInstrumentsIsRunning, 1)
	fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Begin sync instruments")
	s := storage.GetStorage()
	token := s.GetTcsToken()
	if token == "" {
		setLastError(errors.New("No TCS token found"))
		return
	}
	urls := []string{stocksURL, bondsURL, etfURL}
	instruments := []models.Instrument{}
	client := &http.Client{}
	channel := make(chan []models.Instrument)

	for _, url := range urls {
		go getInstruments(client, token, url, channel)
	}

	for range urls {
		instruments = append(instruments, <-channel...)
	}

	_, err := s.DeleteAllInstruments()
	if err != nil {
		setLastError(err)
		return
	}
	err = s.AddInstruments(instruments)
	if err != nil {
		setLastError(err)
		return
	}
	fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Success sync instruments")
}

// GetSyncInstrumentsStatus gets status of instruments sync
func GetSyncInstrumentsStatus() SyncStatus {
	if atomic.LoadInt32(&syncInstrumentsIsRunning) == 1 {
		return SyncStatus{Status: Processing}
	} else if lastSyncIstrumentsError != nil {
		resp := SyncStatus{Status: Err, Error: lastSyncIstrumentsError}
		lastSyncIstrumentsError = nil
		return resp
	} else {
		return SyncStatus{Status: Ok}
	}
}

func setLastError(err error) {
	fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Error sync instruments:", err)
	lastSyncIstrumentsError = err
}

func getInstruments(client *http.Client, token string, url string, c chan []models.Instrument) {
	var respObj struct {
		Payload struct {
			Total int                 `json:"total"`
			Ins   []models.Instrument `json:"instruments"`
		} `json:"payload"`
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c <- nil
		return
	}
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		c <- nil
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c <- nil
		return
	}
	json.Unmarshal(body, &respObj)
	c <- respObj.Payload.Ins
}
