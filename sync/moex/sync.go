package moex

import (
	"errors"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
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
	var prices []priceInt

	from := time.Date(2019, 5, 10, 0, 0, 0, 0, time.UTC)

	cur := 0
	for {
		pr, curRe, err := fetchFromAPI(client, from, ticker, cur)
		if err != nil {
			break
		}
		prices = append(prices, pr...)
		if curRe == 0 {
			break
		}
		cur = curRe
	}
	fmt.Println(len(prices))
}
