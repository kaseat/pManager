package mongo

import (
	"fmt"

	"github.com/kaseat/pManager/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
			"pid":    op.PortfolioID,
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
	res, err := db.operations.InsertMany(ctx, docs, opts)
	if err != nil {
		return err
	}
	if len(res.InsertedIDs) != len(ops) {
		return fmt.Errorf("Not all operations has been inserted")
	}
	return nil
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
