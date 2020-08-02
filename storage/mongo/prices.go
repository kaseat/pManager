package mongo

import (
	"math"
	"time"

	"github.com/kaseat/pManager/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AddPrices saves prices series into a storage
func (db Db) AddPrices(prices []models.Price) error {
	if prices == nil {
		return nil
	}
	if len(prices) == 0 {
		return nil
	}

	docs := make([]interface{}, len(prices))
	for i, p := range prices {
		doc := bson.M{
			"price": int64(math.Round(p.Price * 1e6)),
			"vol":   p.Volume,
			"time":  p.Date,
			"isin":  p.ISIN,
		}
		docs[i] = doc
	}

	ctx := db.context()
	opts := options.InsertMany()
	_, err := db.prices.InsertMany(ctx, docs, opts)
	if err != nil {
		return err
	}

	return nil
}

// GetPrices finds prices depending on input prameters
func (db Db) GetPrices(key string, value string) ([]models.Price, error) {
	filter := bson.M{key: value}
	findOptions := options.Find()
	return db.getPrices(filter, findOptions)
}

// GetPricesByIsin finds prices for given ISIN and dates
func (db Db) GetPricesByIsin(isin, from, to string) ([]models.Price, error) {
	filter := bson.M{"isin": isin}
	and := []interface{}{}
	hasParams := false

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

	return db.getPrices(filter, findOptions)
}

// DeletePrices removes prices depending on input prameters
func (db Db) DeletePrices(key string, value string) (int64, error) {
	filter := bson.M{key: value}
	delOptions := options.Delete()
	return db.delPrices(filter, delOptions)
}

// DeleteAllPrices removes all prices from storage
func (db Db) DeleteAllPrices() (int64, error) {
	filter := bson.M{}
	delOptions := options.Delete()
	return db.delPrices(filter, delOptions)
}

func (db Db) delPrices(filter primitive.M, delOptions *options.DeleteOptions) (int64, error) {
	ctx := db.context()
	del, err := db.prices.DeleteMany(ctx, filter, delOptions)
	if err != nil {
		return 0, err
	}
	return del.DeletedCount, nil
}

func (db Db) getPrices(filter primitive.M, findOptions *options.FindOptions) ([]models.Price, error) {
	ctx := db.context()
	cur, err := db.prices.Find(ctx, filter, findOptions)
	defer cur.Close(ctx)
	if err != nil {
		return nil, err
	}

	var raw []struct {
		Price  int64     `bson:"price"`
		Volume int       `bson:"vol"`
		Date   time.Time `bson:"time"`
		ISIN   string    `bson:"isin"`
	}

	err = cur.All(ctx, &raw)
	if err != nil {
		return nil, err
	}

	if raw == nil {
		return []models.Price{}, nil
	}
	results := make([]models.Price, len(raw))

	for i, it := range raw {
		data := models.Price{
			Price:  float64(it.Price) / 1e6,
			Volume: it.Volume,
			Date:   it.Date,
			ISIN:   it.ISIN,
		}
		results[i] = data
	}

	return results, err
}
