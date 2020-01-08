package tcssync

import (
	"errors"
	"log"
	"time"

	"github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AddOperation adds operation to storage
func AddOperation(op Operation) (string, error) {
	op.OperationID = ""
	res, err := db.operations.InsertOne(db.context(), op)
	if err != nil {
		return "", err
	}

	if id, ok := res.InsertedID.(primitive.ObjectID); ok {
		return id.Hex(), nil
	}

	return "", errors.New("Filed convert 'primitive.ObjectID' to 'string'")
}

// GetOperationByID gets operation by id
func GetOperationByID(operationID string) (Operation, error) {
	var result Operation
	objID, err := primitive.ObjectIDFromHex(operationID)
	if err != nil {
		return result, err
	}
	err = db.operations.FindOne(db.context(), bson.M{"_id": objID}).Decode(&result)
	return result, err
}

// GetAllOperations finds all available operations at the moment
func GetAllOperations() ([]Operation, error) {
	filter := bson.M{}
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"datetime": 1})
	return getOperations(filter, findOptions)
}

// GetAllOperationsByFigi finds all available operations for the specified figi at the moment
func GetAllOperationsByFigi(figi string) ([]Operation, error) {
	filter := bson.M{"figi": figi}
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"datetime": 1})
	return getOperations(filter, findOptions)
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

// GetBalance returns balance of specified currency
func GetBalance(curr Currency) float64 {
	op, _ := GetAllOperations()
	sum := float64(0)
	for _, val := range op {
		if val.Currency != curr {
			continue
		}

		amount := val.Price * float64(val.Quantity)
		if val.OperationType == PayIn {
			sum += amount
		} else {
			sum -= amount
		}
	}
	return sum
}

// SyncPrice updates prices in local storage
func SyncPrice() {
	syncPrices(false)
}

// SyncPriceLastDay updates prices only for the last day
func SyncPriceLastDay() {
	syncPrices(true)
}

func syncPrices(lastDayOnly bool) {
	figis, err := db.operations.Distinct(db.context(), "figi", bson.M{"figi": bson.M{"$ne": RUB}})
	if err != nil {
		log.Fatal("SyncPrice fault: ", err)
		return
	}

	for _, item := range figis {
		if figi, ok := item.(string); ok {
			go updatePriceDb(figi, lastDayOnly)
		}
	}
}

func updatePriceDb(figi string, lastDayOnly bool) {
	var beginDay time.Time
	y, m, d := time.Now().Date()
	endDay := time.Date(y, m, d+1, 0, 0, 0, 0, time.UTC)
	if lastDayOnly {
		beginDay = endDay.AddDate(0, 0, -1)
	} else {
		beginDay = time.Unix(cfg.SyncFrom, 0)
	}

	ctx, cancel := tcs.context()
	candles, err := tcs.client.Candles(ctx, beginDay, endDay, sdk.CandleInterval1Day, figi)
	defer cancel()

	if err != nil {
		log.Fatal("updatePriceDb fault: filed to fetch "+figi+" from tcs API:", err)
		return
	}

	wereErroes := false
	for _, val := range candles {
		filter := bson.M{"$and": []interface{}{
			bson.M{"figi": bson.M{"$eq": val.FIGI}},
			bson.M{"datetime": bson.M{"$eq": val.TS}},
		}}
		update := bson.M{"$set": bson.M{"price": val.ClosePrice, "vol": val.Volume}}
		updateOptions := options.Update()
		updateOptions.SetUpsert(true)
		_, err := db.prices.UpdateOne(db.context(), filter, update, updateOptions)

		if err != nil {
			wereErroes = true
			log.Fatal("updatePriceDb fault: ", err)
		}
	}
	if wereErroes {
		return
	}

	if lastDayOnly {
		log.Println("Sync on", beginDay.Format("2006-01-02"), "for", figi, "complete!")
	} else {
		log.Println("Sync from", beginDay.Format("2006-01-02"), "to", endDay.Format("2006-01-02"), "for", figi, "complete!")
	}
}
