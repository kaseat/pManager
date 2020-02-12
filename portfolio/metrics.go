package portfolio

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetBalance returns blance for specified params
func (p *Portfolio) GetBalance(curr Currency, figi string, on string) (int64, error) {
	pid, err := primitive.ObjectIDFromHex(p.PortfolioID)
	if err != nil {
		return 0, err
	}

	filter := bson.M{"portfolio": pid}
	and := []interface{}{}
	hasParams := false
	if curr != "" {
		and = append(and, bson.M{"currency": curr})
		hasParams = true
	}
	if figi != "" {
		and = append(and, bson.M{"figi": figi})
		hasParams = true
	}

	if curr != "" && figi != "" {
		err = errors.New("You must provide either 'currency' or 'figi'")
		return 0, err
	}

	if dtime, err := time.Parse("2006-01-02T15:04:05.000Z0700", on); err == nil {
		and = append(and, bson.M{"datetime": bson.M{"$lte": dtime}})
		hasParams = true
	}

	if hasParams {
		and = append(and, filter)
		filter = bson.M{"$and": and}
	}

	op, err := getOperations(filter, options.Find())
	if err != nil {
		return 0, err
	}
	sum := getSum(op)
	if figi != "" {
		sum = -sum
	}
	return sum, nil
}

// GetAveragePriceByFigi returns average price of specified figi (FIFO)
func (p *Portfolio) GetAveragePriceByFigi(figi string) (int64, error) {
	pid, err := primitive.ObjectIDFromHex(p.PortfolioID)
	if err != nil {
		return 0, err
	}
	filter := bson.M{"$and": []interface{}{
		bson.M{"portfolio": pid},
		bson.M{"figi": figi},
	}}
	findOptions := options.Find()
	op, err := getOperations(filter, findOptions)
	if err != nil {
		return 0, err
	}
	return getAverage(op), nil
}

// GetAveragePriceByFigiTillDate returns average price of specified figi (FIFO) on specified date
func (p *Portfolio) GetAveragePriceByFigiTillDate(figi string, dt time.Time) (int64, error) {
	pid, err := primitive.ObjectIDFromHex(p.PortfolioID)
	if err != nil {
		return 0, err
	}
	filter := bson.M{"$and": []interface{}{
		bson.M{"portfolio": pid},
		bson.M{"figi": figi},
		bson.M{"datetime": bson.M{"$lte": dt}},
	}}
	findOptions := options.Find()
	op, err := getOperations(filter, findOptions)
	if err != nil {
		return 0, err
	}
	return getAverage(op), nil
}
