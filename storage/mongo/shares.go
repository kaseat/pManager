package mongo

import (
	"sort"
	"time"

	"github.com/kaseat/pManager/models"
	"github.com/kaseat/pManager/models/operation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetShares gets shares
func (db Db) GetShares(pid string, onDate string) ([]models.Share, error) {
	ops, err := db.GetOperations(pid, "", "", "", onDate)
	if err != nil {
		return nil, err
	}
	shares := make(map[string]models.Share)
	for _, op := range ops {
		if i, ok := shares[op.ISIN]; ok {
			if op.OperationType == operation.Buy {
				i.Volume += op.Volume
				shares[op.ISIN] = i
			}
			if op.OperationType == operation.Sell || op.OperationType == operation.Buyback {
				i.Volume -= op.Volume
				shares[op.ISIN] = i
			}
		} else {
			shares[op.ISIN] = models.Share{
				ISIN:   op.ISIN,
				Volume: op.Volume,
			}
		}
	}

	instr, err := db.GetAllInstruments()
	if err != nil {
		return nil, err
	}
	instrMap := make(map[string]models.Instrument)
	for _, ins := range instr {
		instrMap[ins.ISIN] = ins
	}
	result := []models.Share{}
	for _, sh := range shares {
		if sh.ISIN != "" && sh.Volume != 0 {
			sh.Ticker = instrMap[sh.ISIN].Ticker
			result = append(result, sh)
		}
	}

	if onDate == "" {
		onDate = time.Now().Format("2006-01-02T15:04:05Z07:00")
	}

	if dtime, err := time.Parse("2006-01-02T15:04:05Z07:00", onDate); err == nil {
		y, m, d := dtime.Date()
		dt := time.Date(y, m, d, 7, 0, 0, 0, time.UTC)

		and := []interface{}{
			bson.M{"time": bson.M{"$gte": dt.AddDate(0, 0, -10)}},
			bson.M{"time": bson.M{"$lte": dt}},
		}
		filter := bson.M{"$and": and}
		findOptions := options.Find()
		prices, err := db.getPrices(filter, findOptions)
		if err != nil {
			return nil, err
		}

		pricesMap := make(map[string][]models.Price)
		for _, price := range prices {
			if pr, ok := pricesMap[price.ISIN]; ok {
				pr = append(pr, price)
				pricesMap[price.ISIN] = pr
			} else {
				pricesMap[price.ISIN] = []models.Price{price}
			}
		}

		pricesMapLastPrice := make(map[string]models.Price)
		for key, val := range pricesMap {
			sort.Slice(val, func(i, j int) bool { return val[i].Date.After(val[j].Date) })
			pricesMapLastPrice[key] = val[0]
		}

		for i, sh := range result {
			result[i].Date = dt
			result[i].Price = pricesMapLastPrice[sh.ISIN].Price
		}
	}
	return result, nil
}
