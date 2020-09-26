package moex

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/kaseat/pManager/models/currency"
)

type priceInternal struct {
	Currency currency.Type
	Date     time.Time
	Price    float64
	Volume   int
}

func fetchFromAPI(client *http.Client, from time.Time, ticker string, cursor int) ([]priceInternal, int, error) {
	fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Start fetching prices for", ticker, "from", cursor)
	columns := "history.columns=BOARDID,TRADEDATE,LEGALCLOSEPRICE,VOLUME"
	fromStr := fmt.Sprintf("from=%s", from.Format("2006-01-02"))
	start := fmt.Sprintf("start=%d", cursor)
	url := fmt.Sprintf("http://iss.moex.com/iss/history/engines/stock/markets/shares/securities/%s.json?iss.meta=off", ticker)
	url = fmt.Sprintf("%s&%s&%s&%s", url, columns, fromStr, start)

	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil, 0, err
	}
	defer resp.Body.Close()

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, 0, err
	}

	var rawResponse struct {
		History struct {
			Data [][]interface{} `json:"data"`
		} `json:"history"`
		Cursor struct {
			Data [][]float32 `json:"data"`
		} `json:"history.cursor"`
	}

	err = json.Unmarshal(r, &rawResponse)
	if err != nil {
		fmt.Println(err)
		return nil, 0, err
	}

	prices := make([]priceInternal, len(rawResponse.History.Data))

	for i, rawPrice := range rawResponse.History.Data {
		var curr currency.Type
		if rawPrice[0] == "TQTD" ||
			rawPrice[0] == "TQBD" ||
			rawPrice[0] == "EQTD" {
			curr = currency.USD
		} else if rawPrice[0] == "TQBE" ||
			rawPrice[0] == "TQTE" ||
			rawPrice[0] == "EQTU" {
			curr = currency.EUR
		} else {
			curr = currency.RUB
		}

		dtime, _ := time.Parse("2006-01-02", rawPrice[1].(string))

		price := priceInternal{
			Currency: curr,
			Date:     dtime,
			Price:    rawPrice[2].(float64),
			Volume:   int(rawPrice[3].(float64)),
		}
		prices[i] = price
	}
	if int(rawResponse.Cursor.Data[0][1])-cursor-100 > 0 {
		cursor += 100
	} else {
		cursor = 0
	}
	fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Success load", len(prices), "prices for", ticker)
	return prices, cursor, nil
}
