package mongo

import (
	"fmt"
	"log"
	"time"

	"github.com/kaseat/pManager/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

// GetShares gets shares
func (db Db) GetShares(pid string, onDate string) ([]models.Share, error) {
	defer timeTrack(time.Now(), "GetShares")
	p, err := primitive.ObjectIDFromHex(pid)
	if err != nil {
		return nil, fmt.Errorf("Could not decode portfolio Id (%s). Internal error : %s", pid, err)
	}
	dtime := time.Now()
	if t, err := time.Parse("2006-01-02T15:04:05Z07:00", onDate); err == nil {
		dtime = t
	}

	matchStage := bson.D{{Key: "$match", Value: bson.D{
		{Key: "pid", Value: p},
		{Key: "time", Value: bson.M{"$lte": dtime}},
	}}}
	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "isin", Value: true},
		{Key: "vol", Value: bson.M{
			"$switch": bson.D{
				{Key: "branches", Value: []bson.D{
					{{Key: "case", Value: bson.M{"$in": []interface{}{"$type", []string{"buyback", "sell"}}}},
						{Key: "then", Value: bson.M{"$multiply": []interface{}{"$vol", -1}}}},
					{{Key: "case", Value: bson.M{"$eq": []string{"$type", "buy"}}},
						{Key: "then", Value: "$vol"}},
				}},
				{Key: "default", Value: 0},
			},
		}},
		{Key: "price", Value: bson.M{
			"$switch": bson.D{
				{Key: "branches", Value: []bson.D{
					{{Key: "case", Value: bson.M{"$in": []interface{}{"$type", []string{"buyback", "payIn", "accruedInterestSell", "sell"}}}},
						{Key: "then", Value: bson.M{"$multiply": []string{"$vol", "$price"}}}},
					{{Key: "case", Value: bson.M{"$not": bson.M{"$in": []interface{}{"$type", []string{"buyback", "payIn", "accruedInterestSell", "sell"}}}}},
						{Key: "then", Value: bson.M{"$multiply": []interface{}{"$vol", "$price", -1}}}},
				}},
				{Key: "default", Value: 0},
			},
		}},
	}}}
	groupStage := bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: "$isin"},
		{Key: "vol", Value: bson.M{"$sum": "$vol"}},
		{Key: "balance", Value: bson.M{"$sum": "$price"}},
	}}}
	joinPriceStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "prices"},
		{Key: "let", Value: bson.M{"id": "$_id"}},
		{Key: "pipeline", Value: []bson.M{
			{"$match": bson.M{
				"$expr": bson.M{
					"$and": []bson.M{
						{"$eq": []interface{}{"$isin", "$$id"}},
						{"$lte": []interface{}{"$time", dtime}},
					},
				},
			}},
			{"$sort": bson.M{"time": -1}},
			{"$limit": 1},
		}},
		{Key: "as", Value: "priceinfo"},
	}}}
	joinInfoStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "instruments"},
		{Key: "localField", Value: "_id"},
		{Key: "foreignField", Value: "isin"},
		{Key: "as", Value: "securitiesinfo"},
	}}}
	flattenStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "vol", Value: true},
		{Key: "balance", Value: true},
		{Key: "priceinfo", Value: bson.M{"$arrayElemAt": []interface{}{"$priceinfo", 0}}},
		{Key: "securitiesinfo", Value: bson.M{"$arrayElemAt": []interface{}{"$securitiesinfo", 0}}},
	}}}
	finalStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "vol", Value: true},
		{Key: "balance", Value: true},
		{Key: "price", Value: "$priceinfo.price"},
		{Key: "ticker", Value: "$securitiesinfo.ticker"},
		{Key: "name", Value: "$securitiesinfo.name"},
	}}}
	filterZeroVolStage := bson.D{{Key: "$match", Value: bson.M{
		"$expr": bson.M{
			"$or": []bson.M{
				{"$eq": []interface{}{"$_id", "RUB"}},
				{"$ne": []interface{}{"$vol", 0}},
			},
		},
	}}}
	ctx := db.context()
	cur, err := db.operations.Aggregate(ctx, mongo.Pipeline{matchStage, projectStage, groupStage, joinPriceStage, joinInfoStage, flattenStage, finalStage, filterZeroVolStage})
	defer cur.Close(ctx)

	if err != nil {
		return nil, err
	}

	var rawSecurities []struct {
		ISIN    string `bson:"_id"`
		Ticker  string `bson:"ticker"`
		Name    string `bson:"name"`
		Balance int64  `bson:"balance"`
		Price   int64  `bson:"price"`
		Volume  int64  `bson:"vol"`
	}
	err = cur.All(ctx, &rawSecurities)
	if err != nil {
		return nil, err
	}
	result := []models.Share{}

	var balance int64 = 0
	for _, sec := range rawSecurities {
		balance += sec.Balance
		if sec.ISIN != "RUB" {
			result = append(result, models.Share{
				ISIN:   sec.ISIN,
				Ticker: sec.Ticker,
				Price:  float64(sec.Price) / 1e6,
				Volume: sec.Volume,
				Date:   dtime,
			})
		}
	}
	if len(rawSecurities) != 0 {
		result = append(result, models.Share{
			ISIN:   "RUB",
			Ticker: "RUB",
			Price:  float64(balance) / 1e6,
			Volume: 1,
			Date:   dtime,
		})
	}
	return result, nil
}
