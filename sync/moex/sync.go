package moex

import (
	"errors"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/kaseat/pManager/models"
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

	for _, inst := range ins {
		from := inst.PriceUptdTime
		if from.IsZero() {
			from = time.Date(2019, time.May, 1, 0, 0, 0, 0, time.UTC)
		}
		var pricesRaw []priceInternal
		cur := 0
		for {
			pr, curRe, err := fetchFromAPI(client, from, ticker, cur)
			if err != nil {
				break
			}
			pricesRaw = append(pricesRaw, pr...)
			if curRe == 0 {
				break
			}
			cur = curRe
		}
		prices := make([]models.Price, len(pricesRaw))

		for _, priceRaw := range pricesRaw {
			if inst.Currency != priceRaw.Currency {
				break
			}
			price := models.Price{
				Price:  float64(priceRaw.Price),
				Volume: priceRaw.Volume,
				Date:   priceRaw.Date,
				ISIN:   inst.ISIN,
			}
			prices = append(prices, price)
		}
		s.AddPrices(prices)
	}
}
