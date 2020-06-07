package portfolio

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AddOperation adds operation to the portfolio
func (p *Portfolio) AddOperation(op Operation) (string, error) {
	pid, err := primitive.ObjectIDFromHex(p.PortfolioID)
	if err != nil {
		return "", err
	}
	doc := bson.M{
		"portfolio":     pid,
		"currency":      op.Currency,
		"price":         int64(op.PriceF * 1e6),
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

// GetOperations finds operations depending on input prameters
func (p *Portfolio) GetOperations(figi string, from string, to string) ([]Operation, error) {
	pid, err := primitive.ObjectIDFromHex(p.PortfolioID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"portfolio": pid}
	and := []interface{}{}
	hasParams := false
	if figi != "" {
		and = append(and, bson.M{"figi": figi})
		hasParams = true
	}

	if dtime, err := time.Parse("2006-01-02T15:04:05.000Z07:00", from); err == nil {
		and = append(and, bson.M{"datetime": bson.M{"$gte": dtime}})
		hasParams = true
	}

	if dtime, err := time.Parse("2006-01-02T15:04:05.000Z07:00", to); err == nil {
		and = append(and, bson.M{"datetime": bson.M{"$lte": dtime}})
		hasParams = true
	}

	if hasParams {
		and = append(and, filter)
		filter = bson.M{"$and": and}
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.M{"datetime": 1})
	return getOperations(filter, findOptions)
}
