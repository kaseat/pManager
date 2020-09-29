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

type issSecurity struct {
	Boards map[string]currency.Type
	ISIN   string
	Ticker string
	IsBond bool
}

func fetchFromAPI(client *http.Client, from time.Time, sec issSecurity, cursor int) ([]priceInternal, int, error) {
	fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Start fetching prices for", sec.Ticker, "from", cursor)
	var board string
	if sec.IsBond {
		board = "bonds"
	} else {
		board = "shares"
	}
	columns := "history.columns=BOARDID,TRADEDATE,LEGALCLOSEPRICE,VOLUME,FACEVALUE"
	fromStr := fmt.Sprintf("from=%s", from.Format("2006-01-02"))
	start := fmt.Sprintf("start=%d", cursor)
	url := "http://iss.moex.com/iss/history/engines/stock/markets"
	url = fmt.Sprintf("%s/%s/securities/%s.json", url, board, sec.Ticker)
	url = fmt.Sprintf("%s?iss.meta=off&%s&%s&%s", url, columns, fromStr, start)

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
		curr := sec.Boards[rawPrice[0].(string)]
		dtime, _ := time.Parse("2006-01-02", rawPrice[1].(string))
		var priceValue float64
		if rawPrice[2] == nil {
			continue
		}
		if sec.IsBond {
			priceValue = rawPrice[2].(float64) * rawPrice[4].(float64) / 100
		} else {
			priceValue = rawPrice[2].(float64)
		}
		var volValue int
		if rawPrice[3] == nil {
			volValue = 0
		} else {
			volValue = int(rawPrice[3].(float64))
		}

		price := priceInternal{
			Currency: curr,
			Date:     dtime,
			Price:    priceValue,
			Volume:   volValue,
		}
		prices[i] = price
	}
	if int(rawResponse.Cursor.Data[0][1])-cursor-100 > 0 {
		cursor += 100
	} else {
		cursor = 0
	}
	fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Success load", len(prices), "prices for", sec.Ticker)
	return prices, cursor, nil
}

func getSecurityInfo(client *http.Client, ticker string) (issSecurity, error) {
	securitiesURI := "https://iss.moex.com/iss/securities/%s.json?iss.meta=off"
	url := fmt.Sprintf(securitiesURI, ticker)
	var security issSecurity
	response, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
		return security, err
	}
	defer response.Body.Close()

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return security, err
	}

	var rawResponse struct {
		Description struct {
			Data [][]interface{} `json:"data"`
		} `json:"description"`
		Boards struct {
			Data [][]interface{} `json:"data"`
		} `json:"boards"`
	}

	err = json.Unmarshal(responseBytes, &rawResponse)
	if err != nil {
		fmt.Println(err)
		return security, err
	}

	for _, description := range rawResponse.Description.Data {
		if description[0] == "SECID" {
			security.Ticker = description[2].(string)
		}
		if description[0] == "ISIN" {
			security.ISIN = description[2].(string)
		}
		if description[0] == "INITIALFACEVALUE" {
			security.IsBond = true
		}
	}

	security.Boards = map[string]currency.Type{}
	for _, boardInfo := range rawResponse.Boards.Data {
		if security.IsBond {
			if boardInfo[5] == "bonds" && boardInfo[7] == "stock" {
				security.Boards[boardInfo[1].(string)] = currency.Type(boardInfo[15].(string))
			}
		} else {
			if boardInfo[5] == "shares" && boardInfo[7] == "stock" {
				security.Boards[boardInfo[1].(string)] = currency.Type(boardInfo[15].(string))
			}
		}
	}

	return security, nil
}
