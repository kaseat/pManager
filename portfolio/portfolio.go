package portfolio

import (
	"context"
	"errors"
	"time"

	"github.com/oleiade/lane"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Currency represents string currency
type Currency string

// Type is market operation type
type Type string

// Operation represents market operation
type Operation struct {
	PortfolioID   string    `json:"pid,omitempty" bson:"portfolio"`
	OperationID   string    `json:"id" bson:"_id,omitempty"`
	Currency      Currency  `json:"currency"`
	Price         float64   `json:"price"`
	Volume        int64     `json:"vol"`
	FIGI          string    `json:"figi"`
	DateTime      time.Time `json:"date"`
	OperationType Type      `json:"operationType"`
}

// Portfolio represets a range of investments
type Portfolio struct {
	PortfolioID string `json:"id" bson:"_id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Config represents configuration
type Config struct {
	MongoURL string `json:"mongoURL"`
	DbName   string `json:"dbName"`
}

type dbContext func() context.Context

type database struct {
	operations *mongo.Collection
	prices     *mongo.Collection
	portfolios *mongo.Collection
	context    dbContext
}

const (
	// RUB is Russian Rubble
	RUB Currency = "RUB"
	// USD is US Dollar
	USD Currency = "USD"
	// EUR is Euro
	EUR Currency = "EUR"
)

const (
	// Buy operation
	Buy Type = "buy"
	// Sell operation
	Sell Type = "sell"
	// BrokerCommission operation
	BrokerCommission Type = "brokerCommission"
	// ExchangeCommission operation
	ExchangeCommission Type = "exchangeCommission"
	// PayIn operation
	PayIn Type = "payIn"
	// PayOut operation
	PayOut Type = "payOut"
	// Coupon operation
	Coupon Type = "coupon"
)

var db database
var cfg Config

// Init portgolio module
func Init(config Config) error {
	cfg = config
	db.context = func() context.Context { return context.Background() }
	clientOptions := options.Client().ApplyURI(cfg.MongoURL)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return err
	}
	err = client.Connect(db.context())
	if err != nil {
		return err
	}
	db.operations = client.Database(cfg.DbName).Collection("Operations")
	db.portfolios = client.Database(cfg.DbName).Collection("Portfolios")
	db.prices = client.Database(cfg.DbName).Collection("Prices")
	return nil
}

// AddPortfolio adds new potrfolio
func AddPortfolio(name string, description string) (Portfolio, error) {
	portfolio := Portfolio{
		Name:        name,
		Description: description}

	res, err := db.portfolios.InsertOne(db.context(), portfolio)
	if err != nil {
		return portfolio, err
	}

	if id, ok := res.InsertedID.(primitive.ObjectID); ok {
		portfolio.PortfolioID = id.Hex()
		return portfolio, nil
	}

	return portfolio, errors.New("Filed convert 'primitive.ObjectID' to 'string'")
}

// GetPortfolio gets operation by id
func GetPortfolio(portfolioID string) (bool, Portfolio, error) {
	var result Portfolio

	objID, err := primitive.ObjectIDFromHex(portfolioID)
	if err != nil {
		err = errors.New("Invalid portfolio Id")
		return false, result, err
	}

	filter := bson.M{"_id": objID}
	findOptions := options.Find()
	ctx := db.context()

	cur, err := db.portfolios.Find(ctx, filter, findOptions)
	defer cur.Close(ctx)
	if err != nil {
		return false, result, err
	}

	hasResult := cur.TryNext(ctx)
	if !hasResult {
		return false, result, nil
	}

	err = cur.Decode(&result)
	if err != nil {
		return false, result, err
	}
	return true, result, nil
}

// UpdatePortfolio updates current portfolio
func (p *Portfolio) UpdatePortfolio() (bool, error) {

	objID, err := primitive.ObjectIDFromHex(p.PortfolioID)
	if err != nil {
		err = errors.New("Invalid portfolio Id")
		return false, err
	}

	filter := bson.M{"_id": objID}
	update := bson.M{
		"$set": bson.M{
			"name":        p.Name,
			"description": p.Description,
		},
	}
	res, err := db.portfolios.UpdateOne(db.context(), filter, update)
	if err != nil {
		return false, err
	}
	if res.ModifiedCount > 0 {
		return true, nil
	}
	return false, nil
}

// GetAllPortfolios finds all available portfolios at the moment
func GetAllPortfolios() ([]Portfolio, error) {
	filter := bson.M{}
	findOptions := options.Find()

	ctx := db.context()
	cur, err := db.portfolios.Find(ctx, filter, findOptions)
	defer cur.Close(ctx)

	if err != nil {
		return nil, err
	}

	var results []Portfolio

	err = cur.All(ctx, &results)
	if err != nil {
		return nil, err
	}

	return results, err
}

// DeleteAllPortfolios removes all portfolios
func DeleteAllPortfolios() (bool, error) {
	ctx := db.context()
	filter := bson.M{}

	res, err := db.portfolios.DeleteMany(ctx, filter)
	if err != nil {
		return false, err
	}

	if res.DeletedCount != 0 {
		return true, nil
	}

	return false, nil
}

// DeletePortfolio removes portfolio by Id
func DeletePortfolio(portfolioID string) (bool, error) {
	ctx := db.context()

	objID, err := primitive.ObjectIDFromHex(portfolioID)
	if err != nil {
		err = errors.New("Invalid portfolio Id")
		return false, err
	}

	filter := bson.M{"_id": objID}
	res, err := db.portfolios.DeleteOne(ctx, filter)
	if err != nil {
		return false, err
	}

	if res.DeletedCount != 0 {
		return true, err
	}

	return false, err
}

func (p *Portfolio) String() string {
	s := p.PortfolioID + " " + p.Name + " " + p.Description
	return s
}

// AddOperation adds operation to the portfolio
func (p *Portfolio) AddOperation(op Operation) (string, error) {
	pid, err := primitive.ObjectIDFromHex(p.PortfolioID)
	if err != nil {
		return "", err
	}
	doc := bson.M{
		"portfolio":     pid,
		"currency":      op.Currency,
		"price":         op.Price,
		"volume":        op.Volume,
		"figi":          op.FIGI,
		"datetime":      op.DateTime,
		"operationtype": op.OperationType}

	res, err := db.operations.InsertOne(db.context(), doc)
	if err != nil {
		return "", err
	}

	if id, ok := res.InsertedID.(primitive.ObjectID); ok {
		return id.Hex(), nil
	}

	return "", errors.New("Filed convert 'primitive.ObjectID' to 'string'")
}

// DeleteAllOperations removes all operations associated with given portfolio
func (p *Portfolio) DeleteAllOperations() (int64, error) {
	pid, err := primitive.ObjectIDFromHex(p.PortfolioID)
	if err != nil {
		return 0, err
	}
	filter := bson.M{"portfolio": pid}
	res, err := db.operations.DeleteMany(db.context(), filter)
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

// GetOperationByID gets operation by id
func (p *Portfolio) GetOperationByID(operationID string) (Operation, error) {
	var result Operation
	objID, err := primitive.ObjectIDFromHex(operationID)
	if err != nil {
		return result, err
	}
	filter := bson.M{"_id": objID}
	err = db.operations.FindOne(db.context(), filter).Decode(&result)
	return result, err
}

// GetAllOperations finds all available operations at the moment
func (p *Portfolio) GetAllOperations() ([]Operation, error) {
	pid, err := primitive.ObjectIDFromHex(p.PortfolioID)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"portfolio": pid}
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"datetime": 1})
	return getOperations(filter, findOptions)
}

// GetAllOperationsByFigi finds all available operations for the specified figi at the moment
func (p *Portfolio) GetAllOperationsByFigi(figi string) ([]Operation, error) {
	pid, err := primitive.ObjectIDFromHex(p.PortfolioID)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"$and": []interface{}{
		bson.M{"portfolio": pid},
		bson.M{"figi": figi},
	}}
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"datetime": 1})
	return getOperations(filter, findOptions)
}

// GetAllOperationsByFigiTimeBound finds all available operations for the specified figi for specified time range
func (p *Portfolio) GetAllOperationsByFigiTimeBound(figi string, from time.Time, to time.Time) ([]Operation, error) {
	pid, err := primitive.ObjectIDFromHex(p.PortfolioID)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"$and": []interface{}{
		bson.M{"portfolio": pid},
		bson.M{"figi": figi},
		bson.M{"datetime": bson.M{"$gte": from}},
		bson.M{"datetime": bson.M{"$lte": to}},
	}}
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"datetime": 1})
	return getOperations(filter, findOptions)
}

// GetBalanceByCurrency returns balance of specified currency
func (p *Portfolio) GetBalanceByCurrency(curr Currency) (float64, error) {
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
func (p *Portfolio) GetBalanceByCurrencyTillDate(curr Currency, dt time.Time) (float64, error) {
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
func (p *Portfolio) GetBalanceByFigi(figi string) (float64, error) {
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
func (p *Portfolio) GetBalanceByFigiTillDate(figi string, dt time.Time) (float64, error) {
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
func (p *Portfolio) GetAveragePriceByFigi(figi string) (float64, error) {
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
func (p *Portfolio) GetAveragePriceByFigiTillDate(figi string, dt time.Time) (float64, error) {
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

func getSum(operations []Operation) float64 {
	sum := float64(0)
	for _, opertion := range operations {
		amount := opertion.Price * float64(opertion.Volume)
		switch opertion.OperationType {
		case PayIn, Sell:
			sum += amount
		default:
			sum -= amount
		}
	}
	return sum
}

func getOperations(filter primitive.M, findOptions *options.FindOptions) ([]Operation, error) {
	ctx := db.context()
	cur, err := db.operations.Find(ctx, filter, findOptions)
	defer cur.Close(ctx)

	if err != nil {
		return nil, err
	}

	var results []Operation
	cur.All(ctx, &results)
	return results, err
}

func getAverage(ops []Operation) float64 {
	d := lane.NewDeque()
	for _, op := range ops {
		if op.OperationType == Buy {
			d.Append(op)
		} else {
			for {
				if d.Empty() {
					break
				}
				o := d.Shift().(Operation)
				if o.Volume-op.Volume <= 0 {
					op.Volume -= o.Volume
				} else {
					o.Volume -= op.Volume
					d.Prepend(o)
					break
				}
			}
		}
	}

	cost, vol := 0.0, 0.0
	for {
		if d.Empty() {
			break
		}
		op := d.Pop().(Operation)
		v := float64(op.Volume)
		cost += op.Price * v
		vol += v
	}

	result := 0.0
	if vol != 0 {
		result = cost / vol
	}
	return result
}
