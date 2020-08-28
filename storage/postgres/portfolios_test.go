package postgres

import (
	"testing"

	"github.com/kaseat/pManager/models"
)

func TestPortfolioDeletion(t *testing.T) {
	login, email, hash := "login", "email", "hash"
	uid, _ := db.AddUser(login, email, hash)
	p := models.Portfolio{
		Name:        "name",
		Description: "description",
	}

	p1, _ := db.AddPortfolio(uid, p)
	db.AddPortfolio(uid, p)
	db.AddPortfolio(uid, p)

	resBool, err := db.DeletePortfolio(uid, p1)
	if err != nil {
		t.Errorf("Fail! Could not remove test portfolio. Internal error: %s", err)
	}
	if resBool == true {
		t.Logf("Success! Expected %v, got %v", true, resBool)
	} else {
		t.Errorf("Fail! Portfolio did not remove as it should! Expected %v, got %v", true, resBool)
	}

	resArr, err := db.GetPortfolios(uid)
	if err != nil {
		t.Errorf("Fail! Could not get test portfolio. Internal error: %s", err)
	}
	if len(resArr) == 2 {
		t.Logf("Success! Expected %d, got %d", 2, len(resArr))
	} else {
		t.Errorf("Fail! Received number of portfolios not match! Expected %d, got %d", 2, len(resArr))
	}

	resInt, err := db.DeletePortfolios(uid)
	if err != nil {
		t.Errorf("Fail! Could not remove all test portfolios. Internal error: %s", err)
	}
	if resInt == 2 {
		t.Logf("Success! Expected %d, got %d", 2, resInt)
	} else {
		t.Errorf("Fail! All portfolios did not remove as it should! Expected %d, got %d", 2, resInt)
	}

	resArr, err = db.GetPortfolios(uid)
	if err != nil {
		t.Errorf("Fail! Could not get test portfolio. Internal error: %s", err)
	}
	if len(resArr) == 0 {
		t.Logf("Success! Expected %d, got %d", 0, len(resArr))
	} else {
		t.Errorf("Fail! Received number of portfolios not match! Expected %d, got %d", 0, len(resArr))
	}

	malformedID := "ffff"
	// feed DeletePortfolio with malformed portfolio Id
	_, err = db.DeletePortfolio(uid, malformedID)
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

	// feed DeletePortfolio with malformed user Id
	_, err = db.DeletePortfolio(malformedID, p1)
	expectedErrMsg = "Invalid user Id format. Expected positive number"
	if err == nil {
		t.Errorf("Fail! Expected '%s' error", expectedErrMsg)
	} else {
		if err.Error() != expectedErrMsg {
			t.Errorf("Fail! Expected '%s' got '%s'", expectedErrMsg, err)
		} else {
			t.Logf("Success! Expected '%s' got '%s'", expectedErrMsg, err)
		}
	}

	// feed DeletePortfolios with malformed user Id
	resInt, err = db.DeletePortfolios(malformedID)
	if err == nil {
		t.Errorf("Fail! Expected '%s' error", expectedErrMsg)
	} else {
		if err.Error() != expectedErrMsg {
			t.Errorf("Fail! Expected '%s' got '%s'", expectedErrMsg, err)
		} else {
			t.Logf("Success! Expected '%s' got '%s'", expectedErrMsg, err)
		}
	}

	// feed DeletePortfolio with unknown user Id
	unknownID := "0"
	resBool, err = db.DeletePortfolio(unknownID, p1)
	if err != nil {
		t.Errorf("Fail! Could not remove all portfolios. Internal error: %s", err)
	}
	if resBool == false {
		t.Logf("Success! Expected %v, got %v", false, resBool)
	} else {
		t.Errorf("Fail! Portfolio removed, but it should not! Expected %v, got %v", false, resBool)
	}

	// feed DeletePortfolios with unknown user Id
	resInt, err = db.DeletePortfolios(unknownID)
	if err != nil {
		t.Errorf("Fail! Could not remove all portfolios. Internal error: %s", err)
	}
	if resInt == 0 {
		t.Logf("Success! Expected %v, got %v", 0, resInt)
	} else {
		t.Errorf("Fail! Portfolio removed, but it should not! Expected %v, got %v", 0, resInt)
	}

	// feed DeletePortfolio with unknown portfolio Id
	_, err = db.DeletePortfolio(uid, unknownID)
	if err != nil {
		t.Errorf("Fail! Could not remove all portfolios. Internal error: %s", err)
	}
	if resBool == false {
		t.Logf("Success! Expected %v, got %v", false, resBool)
	} else {
		t.Errorf("Fail! Portfolio removed, but it should not! Expected %v, got %v", false, resBool)
	}
	db.DeleteUser(login)
}

func TestPortfolios(t *testing.T) {
	login, email, hash := "login", "email", "hash"
	uid, err := db.AddUser(login, email, hash)
	p := models.Portfolio{
		Name:        "name",
		Description: "description",
	}
	// check adding and getting portfolio
	pid, err := db.AddPortfolio(uid, p)
	if err != nil {
		t.Errorf("Fail! Could not save test portfolio. Internal error: %s", err)
	}

	pid, err = db.AddPortfolio(uid, p)
	if err != nil {
		t.Errorf("Fail! Could not save test portfolio. Internal error: %s", err)
	}

	res, err := db.GetPortfolio(uid, pid)
	if err != nil {
		t.Errorf("Fail! Could not get test portfolio. Internal error: %s", err)
	}
	if res.Name == p.Name && res.Description == p.Description {
		t.Logf("Success! Expected %s, got %s", p, res)
	} else {
		t.Errorf("Fail! Saved and fetched time not match! Expected %s, got %s", p, res)
	}

	// ensure we successfully updated portfolio
	p.Name = "newName"
	resBool, err := db.UpdatePortfolio(uid, pid, p)
	if err != nil {
		t.Errorf("Fail! Could not update test portfolio. Internal error: %s", err)
	}
	if resBool == true {
		t.Logf("Success! Expected %v, got %v", true, resBool)
	} else {
		t.Errorf("Fail! Expected %v, got %v", true, resBool)
	}

	res, err = db.GetPortfolio(uid, pid)
	if err != nil {
		t.Errorf("Fail! Could not get test portfolio. Internal error: %s", err)
	}
	if res.Name == p.Name {
		t.Logf("Success! Expected %s, got %s", p.Name, res.Name)
	} else {
		t.Errorf("Fail! Saved and fetched updated values not match! Expected %s, got %s", p.Name, res.Name)
	}

	// ensure we got all portfolios
	resArr, err := db.GetPortfolios(uid)
	if err != nil {
		t.Errorf("Fail! Could not get test portfolio. Internal error: %s", err)
	}
	if len(resArr) == 2 {
		t.Logf("Success! Expected %d, got %d", 2, len(resArr))
	} else {
		t.Errorf("Fail! Received number of portfolios not match! Expected %d, got %d", 2, len(resArr))
	}

	// feed AddPortfolio with malformed user Id
	malformedID := "ffff"
	_, err = db.AddPortfolio(malformedID, p)
	expectedErrMsg := "Invalid user Id format. Expected positive number"
	if err == nil {
		t.Errorf("Fail! Expected '%s' error", expectedErrMsg)
	} else {
		if err.Error() != expectedErrMsg {
			t.Errorf("Fail! Expected '%s' got '%s'", expectedErrMsg, err)
		} else {
			t.Logf("Success! Expected '%s' got '%s'", expectedErrMsg, err)
		}
	}

	// feed GetPortfolios with malformed user Id
	_, err = db.GetPortfolios(malformedID)
	expectedErrMsg = "Invalid user Id format. Expected positive number"
	if err == nil {
		t.Errorf("Fail! Expected '%s' error", expectedErrMsg)
	} else {
		if err.Error() != expectedErrMsg {
			t.Errorf("Fail! Expected '%s' got '%s'", expectedErrMsg, err)
		} else {
			t.Logf("Success! Expected '%s' got '%s'", expectedErrMsg, err)
		}
	}

	// feed UpdatePortfolio with malformed user Id
	_, err = db.UpdatePortfolio(malformedID, pid, p)
	expectedErrMsg = "Invalid user Id format. Expected positive number"
	if err == nil {
		t.Errorf("Fail! Expected '%s' error", expectedErrMsg)
	} else {
		if err.Error() != expectedErrMsg {
			t.Errorf("Fail! Expected '%s' got '%s'", expectedErrMsg, err)
		} else {
			t.Logf("Success! Expected '%s' got '%s'", expectedErrMsg, err)
		}
	}

	// feed GetPortfolio with malformed portfolio Id
	_, err = db.UpdatePortfolio(uid, malformedID, p)
	expectedErrMsg = "Invalid portfolio Id format. Expected positive number"
	if err == nil {
		t.Errorf("Fail! Expected '%s' error", expectedErrMsg)
	} else {
		if err.Error() != expectedErrMsg {
			t.Errorf("Fail! Expected '%s' got '%s'", expectedErrMsg, err)
		} else {
			t.Logf("Success! Expected '%s' got '%s'", expectedErrMsg, err)
		}
	}

	// feed GetPortfolio with malformed user Id
	_, err = db.GetPortfolio(malformedID, pid)
	expectedErrMsg = "Invalid user Id format. Expected positive number"
	if err == nil {
		t.Errorf("Fail! Expected '%s' error", expectedErrMsg)
	} else {
		if err.Error() != expectedErrMsg {
			t.Errorf("Fail! Expected '%s' got '%s'", expectedErrMsg, err)
		} else {
			t.Logf("Success! Expected '%s' got '%s'", expectedErrMsg, err)
		}
	}

	// feed GetPortfolio with malformed portfolio Id
	_, err = db.GetPortfolio(uid, malformedID)
	expectedErrMsg = "Invalid portfolio Id format. Expected positive number"
	if err == nil {
		t.Errorf("Fail! Expected '%s' error", expectedErrMsg)
	} else {
		if err.Error() != expectedErrMsg {
			t.Errorf("Fail! Expected '%s' got '%s'", expectedErrMsg, err)
		} else {
			t.Logf("Success! Expected '%s' got '%s'", expectedErrMsg, err)
		}
	}

	// feed AddPortfolio with unknown user Id
	unknownID := "0"
	expectedErrMsg = "could not add portfolio to unknown user"
	_, err = db.AddPortfolio(unknownID, p)
	if err == nil {
		t.Errorf("Fail! Expected '%s' error", expectedErrMsg)
	} else {
		if err.Error() != expectedErrMsg {
			t.Errorf("Fail! Expected '%s' got '%s'", expectedErrMsg, err)
		} else {
			t.Logf("Success! Expected '%s' got '%s'", expectedErrMsg, err)
		}
	}

	// feed UpdatePortfolio with unknown user Id
	resBool, err = db.UpdatePortfolio(unknownID, pid, p)
	if err != nil {
		t.Errorf("Fail! Could not update test portfolio. Internal error: %s", err)
	}
	if resBool == false {
		t.Logf("Success! Expected %v, got %v", false, resBool)
	} else {
		t.Errorf("Fail! Expected %v, got %v", false, resBool)
	}

	// feed GetPortfolio with unknown portfolio Id
	_, err = db.UpdatePortfolio(uid, unknownID, p)
	if err != nil {
		t.Errorf("Fail! Could not update test portfolio. Internal error: %s", err)
	}
	if resBool == false {
		t.Logf("Success! Expected %v, got %v", false, resBool)
	} else {
		t.Errorf("Fail! Expected %v, got %v", false, resBool)
	}

	// feed GetPortfolios with unknown user Id
	resArr, err = db.GetPortfolios(unknownID)
	if err != nil {
		t.Errorf("Fail! Could not get all portfolios. Internal error: %s", err)
	}
	if len(resArr) == 0 {
		t.Logf("Success! Expected %d, got %d", 0, len(resArr))
	} else {
		t.Errorf("Fail! Expected %d, got %d", 0, len(resArr))
	}

	// feed GetPortfolio with unknown user Id
	_, err = db.GetPortfolio(unknownID, pid)
	expectedErrMsg = "no rows in result set"
	if err == nil {
		t.Errorf("Fail! Expected '%s' error", err)
	} else {
		if err.Error() != expectedErrMsg {
			t.Errorf("Fail! Expected '%s' got '%s'", expectedErrMsg, err)
		} else {
			t.Logf("Success! Expected '%s' got '%s'", expectedErrMsg, err)
		}
	}

	// feed GetPortfolio with unknown portfolio Id
	_, err = db.GetPortfolio(uid, unknownID)
	if err == nil {
		t.Errorf("Fail! Expected '%s' error", err)
	} else {
		if err.Error() != expectedErrMsg {
			t.Errorf("Fail! Expected '%s' got '%s'", expectedErrMsg, err)
		} else {
			t.Logf("Success! Expected '%s' got '%s'", expectedErrMsg, err)
		}
	}

	// cleanup
	db.DeleteUser(login)
	db.DeletePortfolios(uid)
}
