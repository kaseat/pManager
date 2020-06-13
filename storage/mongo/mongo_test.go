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

func TestLastUpdateTimeStorage(t *testing.T) {
	provider := "test"
	now, _ := time.Parse(time.RFC3339, "2020-05-13T22:08:41Z")
	err := db.SaveLastUpdateTime(provider, now)
	if err != nil {
		t.Errorf("Could not save '%s' provider. Internal error: %s", provider, err)
	}

	res, err := db.GetLastUpdateTime(provider)
	if err != nil {
		t.Errorf("Could not fetch '%s' provider. Internal error: %s", provider, err)
	}
	if res != now {
		t.Errorf("Saved and fetched time not match! Expected %s, got %s", now, res)
	} else {
		t.Logf("Success! Expected %s, got %s", now, res)
	}

	res, err = db.GetLastUpdateTime("unknown")
	if err != nil {
		t.Errorf("Could not fetch 'unknown' provider. Internal error: %s", err)
	}

	if res.IsZero() {
		t.Logf("Success! Expected zero time, got %s", res)
	} else {
		t.Errorf("Error getting unknown provider. Expected zero time, got %s", res)
	}
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
	removeTestOperations(pid)
}

func TestSaveSingleOperation(t *testing.T) {
	pid := addTestPortfolio()
	now, _ := time.Parse(time.RFC3339, "2020-05-13T22:08:41Z")
	ops := getOperations(pid, now)

	err := db.SaveSingleOperation(pid.Hex(), ops[0])
	if err != nil {
		t.Errorf("Fail! Unknown error: '%s'", err)
	}

	res, err := db.GetOperations(pid.Hex(), "", "", "", "")
	if err != nil {
		t.Errorf("Fail! Unknown error: '%s'", err)
	}
	if ops[0].PortfolioID != res[0].PortfolioID {
		t.Errorf("Fail! Expected '%v' got '%v'", ops[0].PortfolioID, res[0].PortfolioID)
	} else {
		t.Logf("Success! Expected '%v' got '%v'", ops[0].PortfolioID, res[0].PortfolioID)
	}

	removeTestPortfolio(pid)
	removeTestOperations(pid)
}

func TestGetOperations(t *testing.T) {
	pid := addTestPortfolio()
	now, _ := time.Parse(time.RFC3339, "2020-05-13T22:08:41Z")
	ops := getOperations(pid, now)

	err := db.SaveMultipleOperations(pid.Hex(), ops)
	if err != nil {
		t.Errorf("Fail! Unknown error: '%s'", err)
	}

	// returns same number of elements as saved
	res, err := db.GetOperations(pid.Hex(), "", "", "", "")
	if err != nil {
		t.Errorf("Fail! Unknown error: '%s'", err)
	}
	if len(res) != len(ops) {
		t.Errorf("Fail! Expected '%d' got '%d'", len(ops), len(res))
	} else {
		t.Logf("Success! Expected '%d' got '%d'", len(ops), len(res))
	}

	// returns filtered by field 'ticker'
	res, err = db.GetOperations(pid.Hex(), "ticker", "FXUS", "", "")
	if err != nil {
		t.Errorf("Fail! Unknown error: '%s'", err)
	}
	if len(res) != 1 {
		t.Errorf("Fail! Expected '1' got '%d'", len(res))
	} else {
		t.Logf("Success! Expected '1' got '%d'", len(res))
	}

	// returns filtered by field 'ticker'
	res, err = db.GetOperations(pid.Hex(), "ticker", "FXUS", "", "")
	if err != nil {
		t.Errorf("Fail! Unknown error: '%s'", err)
	}
	if len(res) != 1 {
		t.Errorf("Fail! Expected '1' got '%d'", len(res))
	} else {
		t.Logf("Success! Expected '1' got '%d'", len(res))
	}

	// check returns empty result when
	res, err = db.GetOperations(pid.Hex(), "ticker", "FXUSS", "", "")
	if err != nil {
		t.Errorf("Fail! Unknown error: '%s'", err)
	}
	if len(res) != 0 {
		t.Errorf("Fail! Expected '0' got '%d'", len(res))
	} else {
		t.Logf("Success! Expected '0' got '%d'", len(res))
	}

	// returns filtered by field 'ticker'
	res, err = db.GetOperations(pid.Hex(), "curr", "RUB", "", "")
	if err != nil {
		t.Errorf("Fail! Unknown error: '%s'", err)
	}
	if len(res) != len(ops) {
		t.Errorf("Fail! Expected '%d' got '%d'", len(ops), len(res))
	} else {
		t.Logf("Success! Expected '%d' got '%d'", len(ops), len(res))
	}

	timeBound := now.AddDate(0, 0, 1).Format("2006-01-02T15:04:05Z")

	// returns operations grater than specified
	res, err = db.GetOperations(pid.Hex(), "", "", timeBound, "")
	if err != nil {
		t.Errorf("Fail! Unknown error: '%s'", err)
	}
	if len(res) != 1 {
		t.Errorf("Fail! Expected '1' got '%v'", timeBound)
	} else {
		t.Logf("Success! Expected '%s'", err)
	}

	// returns operations grater than specified
	res, err = db.GetOperations(pid.Hex(), "", "", "", timeBound)
	if err != nil {
		t.Errorf("Fail! Unknown error: '%s'", err)
	}
	if len(res) != 1 {
		t.Errorf("Fail! Expected '1' got '%v'", timeBound)
	} else {
		t.Logf("Success! Expected '%s'", err)
	}

	// throws error when passing invalid pid
	_, err = db.GetOperations("", "", "", "", "")
	expectedErrMsg := "the provided hex string is not a valid ObjectID"
	if err == nil {
		t.Errorf("Fail! Expected '%s'", expectedErrMsg)
	} else {
		if err.Error() != expectedErrMsg {
			t.Errorf("Fail! Expected '%s' got '%s'", expectedErrMsg, err)
		} else {
			t.Logf("Success! Expected '%s' got '%s'", expectedErrMsg, err)
		}
	}

	// throws error when provided pid is not found
	pidTxt := "5edbc0a72c857652a0542fab"
	_, err = db.GetOperations(pidTxt, "", "", "", "")
	expectedErrMsg = fmt.Sprintf("No portfolio found with %s Id", pidTxt)
	if err == nil {
		t.Errorf("No error! Expected '%s'", expectedErrMsg)
	} else {

		if err.Error() != expectedErrMsg {
			t.Errorf("Fail! Expected '%s' got '%s'", expectedErrMsg, err)
		} else {
			t.Logf("Success! Expected '%s' got '%s'", expectedErrMsg, err)
		}
	}

	removeTestPortfolio(pid)
	removeTestOperations(pid)
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

func removeTestOperations(pid primitive.ObjectID) {
	ctx := db.context()
	filter := bson.M{"pid": pid}
	opts := options.Delete()

	db.portfolios.DeleteMany(ctx, filter, opts)
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
		DateTime:      now.AddDate(0, 0, 2),
		OperationType: models.BrokerageFee,
	}
	ops = append(ops, op1, op2)
	return ops
}
