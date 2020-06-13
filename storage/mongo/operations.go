package mongo

import (
	"fmt"
	"time"

	"github.com/kaseat/pManager/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SaveSingleOperation saves single opertion into a storage
func (db Db) SaveSingleOperation(portfolioID string, op models.Operation) error {
	ops := []models.Operation{op}
	return db.SaveMultipleOperations(portfolioID, ops)
}

// SaveMultipleOperations saves multiple opertions into a storage
func (db Db) SaveMultipleOperations(portfolioID string, ops []models.Operation) error {
	pid, err := db.findPortfolio(portfolioID)
	if err != nil {
		return err
	}
	if pid.IsZero() {
		return fmt.Errorf("No portfolio found with %s Id", portfolioID)
	}

	docs := make([]interface{}, len(ops))

	for i, op := range ops {
		doc := bson.M{
			"pid":    pid,
			"curr":   op.Currency,
			"price":  int64(op.Price * 1e6),
			"vol":    op.Volume,
			"ticker": op.Ticker,
			"time":   op.DateTime,
			"type":   op.OperationType,
		}
		if op.FIGI != "" {
			doc["figi"] = op.FIGI
		}
		if op.ISIN != "" {
			doc["isin"] = op.ISIN
		}
		docs[i] = doc
	}

	ctx := db.context()
	opts := options.InsertMany()
	_, err = db.operations.InsertMany(ctx, docs, opts)
	if err != nil {
		return err
	}
	return nil
}

// GetOperations finds operations depending on input prameters
func (db Db) GetOperations(portfolioID string, key string, value string, from string, to string) ([]models.Operation, error) {
	pid, err := db.findPortfolio(portfolioID)
	if err != nil {
		return []models.Operation{}, err
	}
	if pid.IsZero() {
		return []models.Operation{}, fmt.Errorf("No portfolio found with %s Id", portfolioID)
	}

	filter := bson.M{"pid": pid}
	and := []interface{}{}
	hasParams := false
	if key != "" {
		and = append(and, bson.M{key: value})
		hasParams = true
	}

	if dtime, err := time.Parse("2006-01-02T15:04:05Z07:00", from); err == nil {
		and = append(and, bson.M{"time": bson.M{"$gte": dtime}})
		hasParams = true
	}

	if dtime, err := time.Parse("2006-01-02T15:04:05Z07:00", to); err == nil {
		and = append(and, bson.M{"time": bson.M{"$lte": dtime}})
		hasParams = true
	}

	if hasParams {
		and = append(and, filter)
		filter = bson.M{"$and": and}
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.M{"time": 1})
	return db.getOperations(filter, findOptions)
}

// Checks if portfolio with specified _id exists. Then needs to be checked on .IsZero()
func (db Db) findPortfolio(pid string) (primitive.ObjectID, error) {
	id, err := primitive.ObjectIDFromHex(pid)
	if err != nil {
		return id, err
	}

	ctx := db.context()
	filter := bson.M{"_id": id}
	opts := options.FindOne()

	r := db.portfolios.FindOne(ctx, filter, opts)
	var result struct {
		ID primitive.ObjectID `bson:"_id"`
	}

	if r.Err() != nil {
		return result.ID, nil
	}

	r.Decode(&result)
	return result.ID, nil
}

func (db Db) getOperations(filter primitive.M, findOptions *options.FindOptions) ([]models.Operation, error) {
	ctx := db.context()
	cur, err := db.operations.Find(ctx, filter, findOptions)
	defer cur.Close(ctx)

	if err != nil {
		return nil, err
	}

	var rawOps []struct {
		PortfolioID   string    `bson:"pid"`
		OperationID   string    `bson:"_id"`
		Currency      string    `bson:"curr"`
		Price         int64     `bson:"price"`
		Volume        int64     `bson:"vol"`
		FIGI          string    `bson:"figi"`
		ISIN          string    `bson:"isin"`
		Ticker        string    `bson:"ticker"`
		DateTime      time.Time `bson:"time"`
		OperationType string    `bson:"type"`
	}
	fmt.Println(rawOps)

	cur.All(ctx, &rawOps)

	if rawOps == nil {
		return []models.Operation{}, nil
	}
	results := make([]models.Operation, len(rawOps))

	for i, op := range rawOps {
		data := models.Operation{
			PortfolioID:   op.PortfolioID,
			OperationID:   op.OperationID,
			Currency:      models.Currency(op.Currency),
			Price:         float64(op.Price) / 1e6,
			Volume:        op.Volume,
			FIGI:          op.FIGI,
			ISIN:          op.ISIN,
			Ticker:        op.Ticker,
			DateTime:      op.DateTime,
			OperationType: models.Type(op.OperationType),
		}
		results[i] = data
	}

	return results, err
}
