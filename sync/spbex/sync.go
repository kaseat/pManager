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
)

var isSync int32

// Sync starts spbex sync
func Sync(ticker string, from, to time.Time) {
	defer func() {
		atomic.StoreInt32(&isSync, 0)
		fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "End sync SPBEX")
	}()
	fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Begin sync SPBEX")
	if atomic.LoadInt32(&isSync) == 1 {
		err := errors.New("SPBEX sync already in process")
		fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Error sync instruments:", err)
		return
	}
	atomic.StoreInt32(&isSync, 1)

	if ticker == "" {
		fmt.Println("ticker is empty")
		return
	}

	if from.IsZero() {
		from = time.Now().AddDate(-1, 0, 0)
	}
	if to.IsZero() {
		to = time.Now()
	}

	url := "https://investcab.ru/api/chistory?symbol=%s&resolution=D"
	url = fmt.Sprintf(url, ticker)
	url += fmt.Sprintf("&from=%d&to=%d", from.Unix(), to.Unix())

	fmt.Println(url)

	client := &http.Client{}
	resp, err := client.Get(url)
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
		Price     []float32 `json:"c"`
	}
	err = json.Unmarshal([]byte(strings.Replace(string(r)[1:len(r)-1], `\"`, `"`, -1)), &rawPrices)
	if err != nil {
		fmt.Println(err)
		return
	}

	prices := make([]struct {
		Date  time.Time
		Price float32
	}, len(rawPrices.Price))

	for i := 0; i < len(rawPrices.Price); i++ {
		prices[i] = struct {
			Date  time.Time
			Price float32
		}{
			Date:  time.Unix(int64(rawPrices.Timestamp[i]), 0),
			Price: rawPrices.Price[i],
		}
	}

	fmt.Println(prices)
}
