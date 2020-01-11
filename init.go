package tcssync

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Config represents configuration
type Config struct {
	MongoURL   string `json:"mongoURL"`
	DbName     string `json:"dbName"`
	TcsToken   string `json:"tcsToken"`
	TcsTimeout int32  `json:"tcsTimeout"`
	SyncFrom   int64  `json:"syncFromTimestamp"`
}

type dbContext func() context.Context

type database struct {
	operations *mongo.Collection
	prices     *mongo.Collection
	context    dbContext
}

type tcsContext func() (context.Context, context.CancelFunc)

type tinkoff struct {
	client  *sdk.SandboxRestClient
	context tcsContext
}

var tcs tinkoff
var db database
var cfg Config

// Init filler
func Init(configPath string) {
	loadConfiguration(configPath)
	initDb()
	initTcs()
}

func initTcs() {
	tcs = tinkoff{
		client: sdk.NewSandboxRestClient(cfg.TcsToken),
		context: func() (context.Context, context.CancelFunc) {
			return context.WithTimeout(context.Background(), time.Duration(cfg.TcsTimeout)*time.Second)
		},
	}
}

func initDb() {
	db.context = func() context.Context { return context.Background() }
	clientOptions := options.Client().ApplyURI(cfg.MongoURL)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(db.context())
	if err != nil {
		log.Fatal(err)
	}
	db.operations = client.Database(cfg.DbName).Collection("Operations")
	db.prices = client.Database(cfg.DbName).Collection("Prices")
}

func loadConfiguration(file string) {
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&cfg)
}
