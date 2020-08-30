package postgres

import (
	"testing"
	"time"

	"github.com/kaseat/pManager/models"
	"github.com/kaseat/pManager/models/currency"
	"github.com/kaseat/pManager/models/operation"
)

func TestMultipleOperations(t *testing.T) {
	login, email, hash := "login", "email", "hash"
	uid, err := db.AddUser(login, email, hash)
	pid, err := db.AddPortfolio(uid, models.Portfolio{Name: "bp", Description: "Best Portfolio"})
	now := time.Now()
	ops := getOperations(now)
	addTestSecurities()

	// ensure we can insert multiple operations with no issues
	_, err = db.AddOperations(pid, ops)
	if err != nil {
		t.Errorf("Unknown error: '%s'", err)
	}

	// returns same number of elements as saved
	res, err := db.GetOperations(pid, "", "", "", "")
	if err != nil {
		t.Errorf("Fail! Unknown error: '%s'", err)
	}
	if len(res) != len(ops) {
		t.Errorf("Fail! Expected '%d' got '%d'", len(ops), len(res))
	} else {
		t.Logf("Success! Expected '%d' got '%d'", len(ops), len(res))
	}

	// returns filtered by field 'ticker'
	res, err = db.GetOperations(pid, "ticker", "FXUS", "", "")
	if err != nil {
		t.Errorf("Fail! Unknown error: '%s'", err)
	}
	if len(res) != 1 {
		t.Errorf("Fail! Expected '1' got '%d'", len(res))
	} else {
		t.Logf("Success! Expected '1' got '%d'", len(res))
	}

	// check returns empty result when we provide unknown ticker
	res, err = db.GetOperations(pid, "ticker", "FXGD", "", "")
	if err != nil {
		t.Errorf("Fail! Unknown error: '%s'", err)
	}
	if len(res) != 0 {
		t.Errorf("Fail! Expected '0' got '%d'", len(res))
	} else {
		t.Logf("Success! Expected '0' got '%d'", len(res))
	}

	// returns filtered by field 'figi'
	res, err = db.GetOperations(pid, "figi", "BBG005HLSZ23", "", "")
	if err != nil {
		t.Errorf("Fail! Unknown error: '%s'", err)
	}
	if len(res) != 1 {
		t.Errorf("Fail! Expected '1' got '%d'", len(res))
	} else {
		t.Logf("Success! Expected '1' got '%d'", len(res))
	}

	// check returns empty result when we provide unknown FIGI
	res, err = db.GetOperations(pid, "figi", "BBG0013FFFT4", "", "")
	if err != nil {
		t.Errorf("Fail! Unknown error: '%s'", err)
	}
	if len(res) != 0 {
		t.Errorf("Fail! Expected '0' got '%d'", len(res))
	} else {
		t.Logf("Success! Expected '0' got '%d'", len(res))
	}

	timeBound := now.AddDate(0, 0, 1).Format("2006-01-02T15:04:05Z")

	// returns operations, occurred after provided date
	res, err = db.GetOperations(pid, "", "", timeBound, "")
	if err != nil {
		t.Errorf("Fail! Unknown error: '%s'", err)
	}
	if len(res) != 1 {
		t.Errorf("Fail! Expected '1' got '%d'", len(res))
	} else {
		t.Logf("Success! Expected '1' got '%d'", len(res))
	}

	// returns operations, occurred before provided date
	res, err = db.GetOperations(pid, "", "", "", timeBound)
	if err != nil {
		t.Errorf("Fail! Unknown error: '%s'", err)
	}
	if len(res) != 1 {
		t.Errorf("Fail! Expected '1' got '%d'", len(res))
	} else {
		t.Logf("Success! Expected '1' got '%d'", len(res))
	}

	// feed AddOperations with malformed portfolio Id
	malformedID := "ffff"
	_, err = db.AddOperations(malformedID, ops)
	expectedErrMsg := "Invalid portfolio Id format. Expected positive number"
	if err == nil {
		t.Errorf("Fail! Expected '%s' error", expectedErrMsg)
	} else {
		if err.Error() != expectedErrMsg {
			t.Errorf("Fail! Expected '%s' got '%s'", expectedErrMsg, err)
		} else {
			t.Logf("Success! Expected '%s' got '%s'", expectedErrMsg, err)
		}
	}

	// feed GetOperations with malformed portfolio Id
	_, err = db.GetOperations(malformedID, "", "", "", "")
	if err == nil {
		t.Errorf("Fail! Expected '%s' error", expectedErrMsg)
	} else {
		if err.Error() != expectedErrMsg {
			t.Errorf("Fail! Expected '%s' got '%s'", expectedErrMsg, err)
		} else {
			t.Logf("Success! Expected '%s' got '%s'", expectedErrMsg, err)
		}
	}

	// feed DeleteOperation with malformed portfolio Id
	_, err = db.DeleteOperation(malformedID, res[0].OperationID)
	if err == nil {
		t.Errorf("Fail! Expected '%s' error", expectedErrMsg)
	} else {
		if err.Error() != expectedErrMsg {
			t.Errorf("Fail! Expected '%s' got '%s'", expectedErrMsg, err)
		} else {
			t.Logf("Success! Expected '%s' got '%s'", expectedErrMsg, err)
		}
	}

	// feed DeleteOperations with malformed portfolio Id
	_, err = db.DeleteOperations(malformedID)
	if err == nil {
		t.Errorf("Fail! Expected '%s' error", expectedErrMsg)
	} else {
		if err.Error() != expectedErrMsg {
			t.Errorf("Fail! Expected '%s' got '%s'", expectedErrMsg, err)
		} else {
			t.Logf("Success! Expected '%s' got '%s'", expectedErrMsg, err)
		}
	}

	// feed DeleteOperation with malformed operation Id
	_, err = db.DeleteOperation(pid, malformedID)
	expectedErrMsg = "Invalid operation Id format. Expected positive number"
	if err == nil {
		t.Errorf("Fail! Expected '%s' error", expectedErrMsg)
	} else {
		if err.Error() != expectedErrMsg {
			t.Errorf("Fail! Expected '%s' got '%s'", expectedErrMsg, err)
		} else {
			t.Logf("Success! Expected '%s' got '%s'", expectedErrMsg, err)
		}
	}

	// feed AddOperations with unknwn portfolio Id
	unknownID := "0"
	_, err = db.AddOperations(unknownID, ops)
	expectedErrMsg = "could not add operation to unknown portfolio"
	if err == nil {
		t.Errorf("Fail! Expected '%s'", expectedErrMsg)
	} else {
		if err.Error() != expectedErrMsg {
			t.Errorf("Fail! Expected '%s' got '%s'", expectedErrMsg, err)
		} else {
			t.Logf("Success! Expected '%s' got '%s'", expectedErrMsg, err)
		}
	}

	// throws no error when provided pid is not found
	resInt, err := db.DeleteOperations(unknownID)
	if err != nil {
		t.Errorf("Fail! Unknown error: '%s'", err)
	}
	if resInt != 0 {
		t.Errorf("Fail! Expected '0' got '%d'", resInt)
	} else {
		t.Logf("Success! Expected '0' got '%d'", resInt)
	}

	// throws no error when provided pid is not found
	resBool, err := db.DeleteOperation(unknownID, res[0].OperationID)
	if err != nil {
		t.Errorf("Fail! Unknown error: '%s'", err)
	}
	if resBool != false {
		t.Errorf("Fail! Expected 'false' got '%v'", resBool)
	} else {
		t.Logf("Success! Expected 'false' got '%v'", resBool)
	}

	// throws no error when provided pid is not found
	resBool, err = db.DeleteOperation(pid, unknownID)
	if err != nil {
		t.Errorf("Fail! Unknown error: '%s'", err)
	}
	if resBool != false {
		t.Errorf("Fail! Expected 'false' got '%v'", resBool)
	} else {
		t.Logf("Success! Expected 'false' got '%v'", resBool)
	}

	remOp := res[0].OperationID

	// throws no error when provided pid is not found
	res, err = db.GetOperations(unknownID, "", "", "", "")
	if err != nil {
		t.Errorf("Fail! Unknown error: '%s'", err)
	}
	if len(res) != 0 {
		t.Errorf("Fail! Expected '0' got '%v'", len(res))
	} else {
		t.Logf("Success! Expected '0' got '%v'", len(res))
	}

	// ensure we successfully removed operation with provided Id
	resBool, err = db.DeleteOperation(pid, remOp)
	if err != nil {
		t.Errorf("Fail! Unknown error: '%s'", err)
	}
	if resBool != true {
		t.Errorf("Fail! Expected 'true' got '%v'", resBool)
	} else {
		t.Logf("Success! Expected 'true' got '%v'", resBool)
	}

	// we should get one operation after another one has been deleted
	res, err = db.GetOperations(pid, "", "", "", "")
	if err != nil {
		t.Errorf("Fail! Unknown error: '%s'", err)
	}
	if len(res) != 1 {
		t.Errorf("Fail! Expected '1' got '%v'", len(res))
	} else {
		t.Logf("Success! Expected '1' got '%v'", len(res))
	}

	// ensure we successfully removed all operations
	resInt, err = db.DeleteOperations(pid)
	if err != nil {
		t.Errorf("Fail! Unknown error: '%s'", err)
	}
	if resInt != 1 {
		t.Errorf("Fail! Expected '1' got '%d'", resInt)
	} else {
		t.Logf("Success! Expected '1' got '%d'", resInt)
	}

	// we should get no operations
	res, err = db.GetOperations(pid, "", "", "", "")
	if err != nil {
		t.Errorf("Fail! Unknown error: '%s'", err)
	}
	if len(res) != 0 {
		t.Errorf("Fail! Expected '0' got '%v'", len(res))
	} else {
		t.Logf("Success! Expected '0' got '%v'", len(res))
	}

	db.DeleteUser("login")
	dropTestSecurities()
}

func TestSaveSingleOperation(t *testing.T) {
	login, email, hash := "logins", "emails", "hashs"
	uid, err := db.AddUser(login, email, hash)
	pid, err := db.AddPortfolio(uid, models.Portfolio{Name: "bps", Description: "Best Portfolios"})
	now := time.Now()
	ops := getOperations(now)
	addTestSecurities()

	// ensure we can get back inserted operation
	_, err = db.AddOperation(pid, ops[0])
	if err != nil {
		t.Errorf("Fail! Unknown error: '%s'", err)
	}

	res, err := db.GetOperations(pid, "", "", "", "")
	if err != nil {
		t.Errorf("Fail! Unknown error: '%s'", err)
	}
	if len(res) != 0 {

		if ops[0].ISIN != res[0].ISIN {
			t.Errorf("Fail! Expected '%v' got '%v'", ops[0].ISIN, res[0].ISIN)
		} else {
			t.Logf("Success! Expected '%v' got '%v'", ops[0].ISIN, res[0].ISIN)
		}
	} else {
		t.Errorf("Fail! Expected '%v' got nothing", ops[0].ISIN)
	}

	db.DeleteUser("logins")
	dropTestSecurities()
}

func getOperations(now time.Time) []models.Operation {
	ops := []models.Operation{}
	op1 := models.Operation{
		Currency:      currency.RUB,
		Price:         50.5351,
		Volume:        150,
		ISIN:          "IE00BD3QHZ91",
		Ticker:        "FXUS",
		DateTime:      now,
		OperationType: operation.Buy,
	}

	op2 := models.Operation{
		Currency:      currency.RUB,
		Price:         0.89,
		Volume:        1,
		ISIN:          "RU000A101NZ2",
		Ticker:        "VTBG",
		DateTime:      now.AddDate(0, 0, 2),
		OperationType: operation.Sell,
	}
	ops = append(ops, op1, op2)
	return ops
}

func addTestSecurities() {
	query := "insert into securities (isin, ticker, figi, currency, asset_type, title) values ($1,$2,$3,$4,$5,$6);"
	db.connection.Exec(db.context, query, "IE00BD3QHZ91", "FXUS", "BBG005HLSZ23", "RUB", 31, "FinEx Акции американских компаний")
	query = "insert into securities (isin, ticker, figi, currency, asset_type, title) values ($1,$2,$3,$4,$5,$6);"
	db.connection.Exec(db.context, query, "RU000A101NZ2", "VTBG", "BBG00V9V16J8", "RUB", 35, "ВТБ Фонд Золото")
}

func dropTestSecurities() {
	query := "delete from securities;"
	db.connection.Exec(db.context, query)
}
