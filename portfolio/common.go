package portfolio

import (
	"context"

	"github.com/oleiade/lane"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db database

// Init portgolio module
func Init(config Config) error {
	cfg := config
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
	db.owners = client.Database(cfg.DbName).Collection("Owners")
	db.prices = client.Database(cfg.DbName).Collection("Prices")
	return nil
}

func getSum(operations []Operation) int64 {
	sum := int64(0)
	for _, opertion := range operations {
		amount := opertion.Price * opertion.Volume
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

	if results == nil {
		results = []Operation{}
	}
	return results, err
}

func getAverage(ops []Operation) int64 {
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

	cost, vol := int64(0), int64(0)
	for {
		if d.Empty() {
			break
		}
		op := d.Pop().(Operation)
		cost += op.Price * op.Volume
		vol += op.Volume
	}

	result := int64(0)
	if vol != 0 {
		result = cost / vol
	}
	return result
}
