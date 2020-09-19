package moex

import (
	"errors"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/kaseat/pManager/storage"
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

	s := storage.GetStorage()
	ins, err := s.GetInstruments("code", "MOEX")
	if err != nil {
		fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Error get securities list:", err)
		return
	}
	fmt.Println(ins)
	
	for _, i := range ins {
		from := i.PriceUptdTime
		if from.IsZero() {
			from = time.Date(2019, time.May, 1, 0, 0, 0, 0, time.UTC)
		}
		var prices []priceInternal
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
	}
}
