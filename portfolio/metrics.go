package portfolio

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetBalanceByCurrency returns balance of specified currency
func (p *Portfolio) GetBalanceByCurrency(curr Currency) (int64, error) {
	pid, err := primitive.ObjectIDFromHex(p.PortfolioID)
	if err != nil {
		return 0, err
	}
	filter := bson.M{"$and": []interface{}{
		bson.M{"portfolio": pid},
		bson.M{"currency": curr},
	}}
	findOptions := options.Find()
	op, err := getOperations(filter, findOptions)
	if err != nil {
		return 0, err
	}
	return getSum(op), nil
}

// GetBalanceByCurrencyTillDate returns balance of specified currency till specified date
func (p *Portfolio) GetBalanceByCurrencyTillDate(curr Currency, dt time.Time) (int64, error) {
	pid, err := primitive.ObjectIDFromHex(p.PortfolioID)
	if err != nil {
		return 0, err
	}
	filter := bson.M{"$and": []interface{}{
		bson.M{"portfolio": pid},
		bson.M{"currency": curr},
		bson.M{"datetime": bson.M{"$lte": dt}},
	}}
	op, err := getOperations(filter, options.Find())
	if err != nil {
		return 0, err
	}
	return getSum(op), nil
}

// GetBalanceByFigi returns balance of specified figi
func (p *Portfolio) GetBalanceByFigi(figi string) (int64, error) {
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
	return -getSum(op), nil
}

// GetBalanceByFigiTillDate returns balance of specified figi till specified date
func (p *Portfolio) GetBalanceByFigiTillDate(figi string, dt time.Time) (int64, error) {
	pid, err := primitive.ObjectIDFromHex(p.PortfolioID)
	if err != nil {
		return 0, err
	}
	filter := bson.M{"$and": []interface{}{
		bson.M{"portfolio": pid},
		bson.M{"figi": bson.M{"$eq": figi}},
		bson.M{"datetime": bson.M{"$lte": dt}},
	}}
	findOptions := options.Find()
	op, err := getOperations(filter, findOptions)
	if err != nil {
		return 0, err
	}
	return -getSum(op), nil
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
