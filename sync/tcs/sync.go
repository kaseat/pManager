package sync

import (
	"context"
	"log"
	"time"

	"github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type tcsContext func() (context.Context, context.CancelFunc)

type tinkoff struct {
	client  *sdk.SandboxRestClient
	context tcsContext
}

type dbContext func() context.Context

type database struct {
	operations *mongo.Collection
	prices     *mongo.Collection
	context    dbContext
}

// Config represents configuration
type Config struct {
	MongoURL   string `json:"mongoURL"`
	DbName     string `json:"dbName"`
	TcsToken   string `json:"tcsToken"`
	TcsTimeout int32  `json:"tcsTimeout"`
}

var tcs tinkoff
var db database
var cfg Config

// Init sync module
func Init(config Config) error {
	cfg = config

	// init db
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
	db.operations = client.Database(cfg.DbName).Collection("operations")
	db.prices = client.Database(cfg.DbName).Collection("prices")

	// init tcs
	tcs = tinkoff{
		client: sdk.NewSandboxRestClient(cfg.TcsToken),
		context: func() (context.Context, context.CancelFunc) {
			return context.WithTimeout(context.Background(), time.Duration(cfg.TcsTimeout)*time.Second)
		},
	}
	return nil
}

// Price updates prices in local storage
func Price() {
	syncPrices(false)
}

// PriceLastDay updates prices only for the last day
func PriceLastDay() {
	syncPrices(true)
}

func syncPrices(lastDayOnly bool) {
	filter := bson.M{"figi": bson.M{"$ne": "RUB"}}
	figis, err := db.operations.Distinct(db.context(), "figi", filter)
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
		beginDay = time.Date(2019, 12, 10, 0, 0, 0, 0, time.UTC)
	}

	ctx, cancel := tcs.context()
	candles, err := tcs.client.Candles(ctx, beginDay, endDay, sdk.CandleInterval1Day, figi)
	defer cancel()

	if err != nil {
		log.Fatal("updatePriceDb fault: filed to fetch", figi, "from tcs API:", err)
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
