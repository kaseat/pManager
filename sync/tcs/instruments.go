package tcs

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/kaseat/pManager/models"
	"github.com/kaseat/pManager/storage"
)

const stocksURL = "https://api-invest.tinkoff.ru/openapi/sandbox/market/stocks"
const bondsURL = "https://api-invest.tinkoff.ru/openapi/sandbox/market/bonds"
const etfURL = "https://api-invest.tinkoff.ru/openapi/sandbox/market/etfs"

// SyncInstruments start sync instruments from tcs API
func SyncInstruments() error {
	s := storage.GetStorage()
	token := s.GetTcsToken()
	if token == "" {
		return errors.New("No TCS token found")
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
		return err
	}
	err = s.AddInstruments(instruments)
	if err != nil {
		return err
	}
	return nil
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
