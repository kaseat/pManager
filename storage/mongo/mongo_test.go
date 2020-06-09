package mongo

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/kaseat/pManager/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db Db

func TestMain(m *testing.M) {
	db = Db{}
	db.Init(Config{
		MongoURL: "mongodb://localhost:27017",
		DbName:   "pm_test",
	})

	os.Exit(m.Run())
}

func TestSaveLastUpdateTime(t *testing.T) {
	provider := "test"
	now, _ := time.Parse(time.RFC3339, "2020-05-13T22:08:41Z")
	err := db.SaveLastUpdateTime(provider, now)
	if err != nil {
		t.Errorf("Could not save '%s' provider. Internl error: %s", provider, err)
	}
	res, err := db.GetLastUpdateTime(provider)
	if err != nil {
		t.Errorf("Could not fetch '%s' provider. Internl error: %s", provider, err)
	}
	if res != now {
		t.Errorf("Saved and fetched time not match! Expected %s, got %s", now, res)
	} else {
		t.Logf("Success! Expected %s, got %s", now, res)
	}
}

func TestFindPortfolio(t *testing.T) {

	_, err := db.findPortfolio("")

	expectedErrMsg := "the provided hex string is not a valid ObjectID"
	if err == nil {
		t.Errorf("No error! Expected '%s'", expectedErrMsg)
	} else {

		if err.Error() != expectedErrMsg {
			t.Errorf("Expected '%s' got '%s'", expectedErrMsg, err)
		} else {
			t.Logf("Success! Expected '%s'", err)
		}
	}

	id, err := db.findPortfolio("5edbc0a72c857652a0542fab")

	if err != nil {
		t.Errorf("Got unexpected error '%s'", err)
	} else {
		if id.IsZero() {
			t.Log("Success! Got zero ObjectID as expected")
		} else {
			t.Errorf("Expected zero ObjectID, got '%s'", id.String())
		}
	}

	pid := addTestPortfolio()

	id, err = db.findPortfolio(pid.Hex())
	if err != nil {
		t.Errorf("Got unexpected error '%s'", err)
	} else {
		if id == pid {
			t.Logf("Success! Expected '%s', got '%s'", pid.String(), id.String())
		} else {
			t.Errorf("Found portfolio Id not match! Expected '%s', got '%s'", pid.String(), id.String())
		}
	}

	removeTestPortfolio(pid)
}

func TestSaveMultipleOperations(t *testing.T) {
	ops := []models.Operation{}

	err := db.SaveMultipleOperations("", ops)
	expectedErrMsg := "the provided hex string is not a valid ObjectID"
	if err == nil {
		t.Errorf("No error! Expected '%s'", expectedErrMsg)
	} else {

		if err.Error() != expectedErrMsg {
			t.Errorf("Expected '%s' got '%s'", expectedErrMsg, err)
		} else {
			t.Logf("Success! Expected '%s'", err)
		}
	}

	pidTxt := "5edbc0a72c857652a0542fab"
	err = db.SaveMultipleOperations(pidTxt, ops)
	expectedErrMsg = fmt.Sprintf("No portfolio found with %s Id", pidTxt)
	if err == nil {
		t.Errorf("No error! Expected '%s'", expectedErrMsg)
	} else {

		if err.Error() != expectedErrMsg {
			t.Errorf("Expected '%s' got '%s'", expectedErrMsg, err)
		} else {
			t.Logf("Success! Expected '%s'", err)
		}
	}
	pid := addTestPortfolio()
	now, _ := time.Parse(time.RFC3339, "2020-05-13T22:08:41Z")
	ops = getOperations(pid, now)

	err = db.SaveMultipleOperations(pid.Hex(), ops)

	if err != nil {
		t.Errorf("Unknown error: '%s'", err)
	}

	removeTestPortfolio(pid)
	db.operations.DeleteMany(db.context(), bson.M{"pid": pid}, options.Delete())
}

func TestGetOperations(t *testing.T) {
	pid := addTestPortfolio()
	now, _ := time.Parse(time.RFC3339, "2020-05-13T22:08:41Z")
	ops := getOperations(pid, now)

	err := db.SaveMultipleOperations(pid.Hex(), ops)
	if err != nil {
		t.Errorf("Unknown error: '%s'", err)
	}

	res, err := db.GetOperations(pid.Hex(), "", "", "", "")
	if err != nil {
		t.Errorf("Unknown error: '%s'", err)
	}

	if len(res) != len(ops) {
		t.Errorf("Expected '%d' got '%d'", len(ops), len(res))
	} else {
		t.Logf("Success! Expected '%s'", err)
	}

	removeTestPortfolio(pid)
	db.operations.DeleteMany(db.context(), bson.M{"pid": pid}, options.Delete())
}

func addTestPortfolio() primitive.ObjectID {
	ctx := db.context()
	opts := options.InsertOne()
	var testItem interface{} = struct {
		TestKey string `bson:"test_key"`
	}{"test_value"}

	r, _ := db.portfolios.InsertOne(ctx, testItem, opts)
	return r.InsertedID.(primitive.ObjectID)
}

func removeTestPortfolio(id primitive.ObjectID) {
	ctx := db.context()
	filter := bson.M{"_id": id}
	opts := options.Delete()

	db.portfolios.DeleteOne(ctx, filter, opts)
}

func getOperations(pid primitive.ObjectID, now time.Time) []models.Operation {
	ops := []models.Operation{}
	op1 := models.Operation{
		PortfolioID:   pid.Hex(),
		Currency:      models.RUB,
		Price:         50.5351,
		Volume:        150,
		ISIN:          "IE00BD3QHZ91",
		Ticker:        "FXUS",
		DateTime:      now,
		OperationType: models.Buy,
	}

	op2 := models.Operation{
		PortfolioID:   pid.Hex(),
		Currency:      models.RUB,
		Price:         0.89,
		Volume:        1,
		FIGI:          "BBG0013HGFT4",
		Ticker:        "RUB",
		DateTime:      now,
		OperationType: models.BrokerageFee,
	}
	ops = append(ops, op1, op2)
	return ops
}
