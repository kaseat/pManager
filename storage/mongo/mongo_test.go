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
	"golang.org/x/oauth2"
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

func TestUsers(t *testing.T) {
	login, email, hash := "login", "email", "hash"

	uid, err := db.AddUser(login, email, hash)
	if err != nil {
		t.Errorf("Fail! Could not add test user. Internal error: %s", err)
	}
	pass, err := db.GetUserPassword(login)
	if err != nil {
		t.Errorf("Fail! Could not fetsh test user's password. Internal error: %s", err)
	}
	if pass == hash {
		t.Logf("Success! Expected %v, got %v", hash, pass)
	} else {
		t.Errorf("Fail! Saved and fetched passwords not match! Expected %v, got %v", hash, pass)
	}
	user, err := db.GetUserByLogin(login)
	if err != nil {
		t.Errorf("Fail! Could not get user by login. Internal error: %s", err)
	}
	if user.UserID == uid {
		t.Logf("Success! Expected %v, got %v", uid, user.UserID)
	} else {
		t.Errorf("Fail! Saved and fetched user Ids not match! Expected %v, got %v", uid, user.UserID)
	}

	state := "state"
	err = db.AddUserState(login, state)
	if err != nil {
		t.Errorf("Fail! Could not add user state. Internal error: %s", err)
	}
	st, err := db.GetUserState(login)
	if err != nil {
		t.Errorf("Fail! Could not get user's state. Internal error: %s", err)
	}
	if state == st {
		t.Logf("Success! Expected %v, got %v", state, st)
	} else {
		t.Errorf("Fail! Saved and fetched user states not match! Expected %v, got %v", state, st)
	}

	exp := oauth2.Token{}
	tok, err := db.GetUserToken(login)
	if err != nil {
		t.Errorf("Fail! Could not get user's token. Internal error: %s", err)
	}
	if tok == exp {
		t.Logf("Success! Expected %v, got %v", exp, tok)
	} else {
		t.Errorf("Fail! Saved and fetched user tokens not match! Expected %v, got %v", exp, tok)
	}

	token := oauth2.Token{
		AccessToken:  "access_token",
		TokenType:    "Bearer",
		RefreshToken: "refresh_token",
		Expiry:       time.Now().Round(5),
	}
	err = db.AddUserToken(state, &token)
	if err != nil {
		t.Errorf("Fail! Could not add user token. Internal error: %s", err)
	}
	tok, err = db.GetUserToken(login)
	if err != nil {
		t.Errorf("Fail! Could not get user's token. Internal error: %s", err)
	}
	if token == tok {
		t.Logf("Success! Expected %v, got %v", token, tok)
	} else {
		t.Errorf("Fail! Saved and fetched user tokens not match! Expected %v, got %v", token, tok)
	}

	_, err = db.AddUser(login, "", hash)
	expectedErrMsg := "User with this login already exists"
	if err == nil {
		t.Errorf("Fail! Expected '%s' error", expectedErrMsg)
	} else {
		if err.Error() != expectedErrMsg {
			t.Errorf("Fail! Expected '%s' got '%s'", expectedErrMsg, err)
		} else {
			t.Logf("Success! Expected '%s' got '%s'", expectedErrMsg, err)
		}
	}

	newHash := "newHash"
	hasUpdated, err := db.UpdateUserPassword(login, newHash)
	if err != nil {
		t.Errorf("Fail! Could not update test user's paassword. Internal error: %s", err)
	}
	if hasUpdated {
		t.Logf("Success! Expected %v, got %v", true, hasUpdated)
	} else {
		t.Errorf("Fail! Did not update test user's paassword! Expected %v, got %v", true, hasUpdated)
	}

	pass, err = db.GetUserPassword(login)
	if err != nil {
		t.Errorf("Fail! Could not fetsh test user's password after update. Internal error: %s", err)
	}
	if pass == newHash {
		t.Logf("Success! Expected %v, got %v", newHash, pass)
	} else {
		t.Errorf("Fail! Saved and fetched passwords not match after update! Expected %v, got %v", newHash, pass)
	}
	newUser := models.User{
		Login: "newLogin",
		Email: "NewEmail",
	}

	hasUpdated, err = db.UpdateUser(login, newUser)
	if err != nil {
		t.Errorf("Fail! Could not update test user. Internal error: %s", err)
	}
	if hasUpdated {
		t.Logf("Success! Expected %v, got %v", true, hasUpdated)
	} else {
		t.Errorf("Fail! Did not update test user! Expected %v, got %v", true, hasUpdated)
	}

	hasUpdated, err = db.UpdateUser(login, newUser)
	if err != nil {
		t.Errorf("Fail! Could not update test user. Internal error: %s", err)
	}
	if !hasUpdated {
		t.Logf("Success! Expected %v, got %v", false, hasUpdated)
	} else {
		t.Errorf("Fail! Did update test user! Expected %v, got %v", false, hasUpdated)
	}

	hasUpdated, err = db.UpdateUserPassword(login, newHash)
	if err != nil {
		t.Errorf("Fail! Could not update test user. Internal error: %s", err)
	}
	if !hasUpdated {
		t.Logf("Success! Expected %v, got %v", false, hasUpdated)
	} else {
		t.Errorf("Fail! Did update test user's password! Expected %v, got %v", false, hasUpdated)
	}

	pass, err = db.GetUserPassword(login)
	if err != nil {
		t.Errorf("Fail! Could not fetsh test user's password. Internal error: %s", err)
	}
	if pass == "" {
		t.Logf("Success! Expected %v, got %v", "", pass)
	} else {
		t.Errorf("Fail! Saved and fetched passwords not match! Expected %v, got %v", "", pass)
	}

	hasDeleted, err := db.DeleteUser(newUser.Login)
	if err != nil {
		t.Errorf("Fail! Could not remove test user. Internal error: %s", err)
	}
	if hasDeleted {
		t.Logf("Success! Expected %v, got %v", true, hasDeleted)
	} else {
		t.Errorf("Fail! Did not remove test user! Expected %v, got %v", true, hasDeleted)
	}
}

func TestPortfolioDeletion(t *testing.T) {
	// arrange
	u := addTestUser()
	p := models.Portfolio{
		Name:        "name",
		Description: "description",
	}

	p1, _ := db.AddPortfolio(u.Hex(), p)
	p2, _ := db.AddPortfolio(u.Hex(), p)
	p3, _ := db.AddPortfolio(u.Hex(), p)

	ops1 := getOperations(getTime())
	ops2 := getOperations(getTime())
	ops3 := getOperations(getTime())

	db.AddOperations(p1, ops1)
	db.AddOperations(p2, ops2)
	db.AddOperations(p3, ops3)

	//
	resBool, err := db.DeletePortfolio(u.Hex(), p1)
	if err != nil {
		t.Errorf("Fail! Could not remove test portfolio. Internal error: %s", err)
	}
	if resBool == true {
		t.Logf("Success! Expected %v, got %v", true, resBool)
	} else {
		t.Errorf("Fail! Portfolio did not remove as it should! Expected %v, got %v", true, resBool)
	}

	resArr, err := db.GetPortfolios(u.Hex())
	if err != nil {
		t.Errorf("Fail! Could not get test portfolio. Internal error: %s", err)
	}
	if len(resArr) == 2 {
		t.Logf("Success! Expected %d, got %d", 2, len(resArr))
	} else {
		t.Errorf("Fail! Received number of portfolios not match! Expected %d, got %d", 2, len(resArr))
	}

	resInt, err := db.DeletePortfolios(u.Hex())
	if err != nil {
		t.Errorf("Fail! Could not remove all test portfolios. Internal error: %s", err)
	}
	if resInt == 2 {
		t.Logf("Success! Expected %d, got %d", 2, resInt)
	} else {
		t.Errorf("Fail! All portfolios did not remove as it should! Expected %d, got %d", 2, resInt)
	}

	resArr, err = db.GetPortfolios(u.Hex())
	if err != nil {
		t.Errorf("Fail! Could not get test portfolio. Internal error: %s", err)
	}
	if len(resArr) == 0 {
		t.Logf("Success! Expected %d, got %d", 0, len(resArr))
	} else {
		t.Errorf("Fail! Received number of portfolios not match! Expected %d, got %d", 0, len(resArr))
	}

	ops, _ := db.GetOperations(p2, "", "", "", "")
	if len(ops) == 0 {
		t.Logf("Success! Expected %d, got %d", 2, resInt)
	} else {
		t.Errorf("Fail! All portfolios did not remove as it should! Expected %d, got %d", 2, resInt)
	}

	// feed DeletePortfolio with malformed user Id
	malformedID := "ffff"
	_, err = db.DeletePortfolio(malformedID, p1)
	expectedErrMsg := fmt.Sprintf("Could not decode user Id (%s). Internal error : the provided hex string is not a valid ObjectID", malformedID)
	if err == nil {
		t.Errorf("Fail! Expected '%s' error", expectedErrMsg)
	} else {
		if err.Error() != expectedErrMsg {
			t.Errorf("Fail! Expected '%s' got '%s'", expectedErrMsg, err)
		} else {
			t.Logf("Success! Expected '%s' got '%s'", expectedErrMsg, err)
		}
	}

	// feed DeletePortfolio with malformed portfolio Id
	_, err = db.DeletePortfolio(u.Hex(), malformedID)
	expectedErrMsg = fmt.Sprintf("Could not decode portfolio Id (%s). Internal error : the provided hex string is not a valid ObjectID", malformedID)
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
	unknownID := "5edbc0a72c857652a0542fab"
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
	_, err = db.DeletePortfolio(u.Hex(), unknownID)
	expectedErrMsg = "mongo: no documents in result"
	if err != nil {
		t.Errorf("Fail! Could not remove all portfolios. Internal error: %s", err)
	}
	if resBool == false {
		t.Logf("Success! Expected %v, got %v", false, resBool)
	} else {
		t.Errorf("Fail! Portfolio removed, but it should not! Expected %v, got %v", false, resBool)
	}

}

func TestPortfolios(t *testing.T) {
	// arrange
	u := addTestUser()
	p := models.Portfolio{
		Name:        "name",
		Description: "description",
	}
	// check adding and getting portfolio
	pid, err := db.AddPortfolio(u.Hex(), p)
	if err != nil {
		t.Errorf("Fail! Could not save test portfolio. Internal error: %s", err)
	}

	pid, err = db.AddPortfolio(u.Hex(), p)
	if err != nil {
		t.Errorf("Fail! Could not save test portfolio. Internal error: %s", err)
	}

	res, err := db.GetPortfolio(u.Hex(), pid)
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
	resBool, err := db.UpdatePortfolio(u.Hex(), pid, p)
	if err != nil {
		t.Errorf("Fail! Could not update test portfolio. Internal error: %s", err)
	}
	if resBool == true {
		t.Logf("Success! Expected %v, got %v", true, resBool)
	} else {
		t.Errorf("Fail! Expected %v, got %v", true, resBool)
	}

	res, err = db.GetPortfolio(u.Hex(), pid)
	if err != nil {
		t.Errorf("Fail! Could not get test portfolio. Internal error: %s", err)
	}
	if res.Name == p.Name {
		t.Logf("Success! Expected %s, got %s", p.Name, res.Name)
	} else {
		t.Errorf("Fail! Saved and fetched updated values not match! Expected %s, got %s", p.Name, res.Name)
	}

	// ensure we got all portfolios
	resArr, err := db.GetPortfolios(u.Hex())
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
	expectedErrMsg := "the provided hex string is not a valid ObjectID"
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
	expectedErrMsg = fmt.Sprintf("Could not decode user Id (%s). Internal error : the provided hex string is not a valid ObjectID", malformedID)
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
	expectedErrMsg = fmt.Sprintf("Could not decode user Id (%s). Internal error : the provided hex string is not a valid ObjectID", malformedID)
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
	_, err = db.UpdatePortfolio(u.Hex(), malformedID, p)
	expectedErrMsg = fmt.Sprintf("Could not decode portfolio Id (%s). Internal error : the provided hex string is not a valid ObjectID", malformedID)
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
	expectedErrMsg = fmt.Sprintf("Could not decode user Id (%s). Internal error : the provided hex string is not a valid ObjectID", malformedID)
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
	_, err = db.GetPortfolio(u.Hex(), malformedID)
	expectedErrMsg = fmt.Sprintf("Could not decode portfolio Id (%s). Internal error : the provided hex string is not a valid ObjectID", malformedID)
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
	unknownID := "5edbc0a72c857652a0542fab"
	expectedErrMsg = fmt.Sprintf("No user found with %s Id", unknownID)
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
	_, err = db.UpdatePortfolio(u.Hex(), unknownID, p)
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
	expectedErrMsg = "mongo: no documents in result"
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
	_, err = db.GetPortfolio(u.Hex(), unknownID)
	expectedErrMsg = "mongo: no documents in result"
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
	removeTestUser(u)
	removeTestPortfolios(u)
}

func TestLastUpdateTimeStorage(t *testing.T) {
	// arrange
	provider := "test"
	now, _ := time.Parse(time.RFC3339, "2020-05-13T22:08:41Z")

	// ensure we can get back inserted lastUpdateTime
	err := db.AddLastUpdateTime(provider, now)
	if err != nil {
		t.Errorf("Fail! Could not save '%s' provider. Internal error: %s", provider, err)
	}

	res, err := db.GetLastUpdateTime(provider)
	if err != nil {
		t.Errorf("Fail! Could not fetch '%s' provider. Internal error: %s", provider, err)
	}
	if res != now {
		t.Errorf("Fail! Saved and fetched time not match! Expected %s, got %s", now, res)
	} else {
		t.Logf("Success! Expected %s, got %s", now, res)
	}

	// check if we actually deleted lastUpdateTime entry
	err = db.DeleteLastUpdateTime(provider)
	if err != nil {
		t.Errorf("Fail! Could not delete '%s' provider. Internal error: %s", provider, err)
	}

	res, err = db.GetLastUpdateTime(provider)
	if err != nil {
		t.Errorf("Fail! Could not fetch '%s' provider. Internal error: %s", provider, err)
	}

	if res.IsZero() {
		t.Logf("Success! Expected zero time, got %s", res)
	} else {
		t.Errorf("Fail! Error getting '%s' provider. Expected zero time, got %s", provider, err)
	}
}

func TestMultipleOperations(t *testing.T) {
	pid := addTestPortfolio()
	now := getTime()
	ops := getOperations(now)

	// ensure we can insert multiple operations with no issues
	_, err := db.AddOperations(pid.Hex(), ops)
	if err != nil {
		t.Errorf("Unknown error: '%s'", err)
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

	// check returns empty result when we provide unknown ticker
	res, err = db.GetOperations(pid.Hex(), "ticker", "FXGD", "", "")
	if err != nil {
		t.Errorf("Fail! Unknown error: '%s'", err)
	}
	if len(res) != 0 {
		t.Errorf("Fail! Expected '0' got '%d'", len(res))
	} else {
		t.Logf("Success! Expected '0' got '%d'", len(res))
	}

	// returns filtered by field 'figi'
	res, err = db.GetOperations(pid.Hex(), "figi", "BBG0013HGFT4", "", "")
	if err != nil {
		t.Errorf("Fail! Unknown error: '%s'", err)
	}
	if len(res) != 1 {
		t.Errorf("Fail! Expected '1' got '%d'", len(res))
	} else {
		t.Logf("Success! Expected '1' got '%d'", len(res))
	}

	// check returns empty result when we provide unknown FIGI
	res, err = db.GetOperations(pid.Hex(), "figi", "BBG0013FFFT4", "", "")
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
	res, err = db.GetOperations(pid.Hex(), "", "", timeBound, "")
	if err != nil {
		t.Errorf("Fail! Unknown error: '%s'", err)
	}
	if len(res) != 1 {
		t.Errorf("Fail! Expected '1' got '%d'", len(res))
	} else {
		t.Logf("Success! Expected '1' got '%d'", len(res))
	}

	// returns operations, occurred before provided date
	res, err = db.GetOperations(pid.Hex(), "", "", "", timeBound)
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
	expectedErrMsg := "the provided hex string is not a valid ObjectID"
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
	expectedErrMsg = fmt.Sprintf("Could not decode portfolio Id (%s). Internal error : the provided hex string is not a valid ObjectID", malformedID)
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
	expectedErrMsg = fmt.Sprintf("Could not decode portfolio Id (%s). Internal error : the provided hex string is not a valid ObjectID", malformedID)
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
	_, err = db.DeleteOperation(pid.Hex(), malformedID)
	expectedErrMsg = fmt.Sprintf("Could not decode operation Id (%s). Internal error : the provided hex string is not a valid ObjectID", malformedID)
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
	expectedErrMsg = fmt.Sprintf("Could not decode portfolio Id (%s). Internal error : the provided hex string is not a valid ObjectID", malformedID)
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
	unknownID := "5edbc0a72c857652a0542fab"
	_, err = db.AddOperations(unknownID, ops)
	expectedErrMsg = fmt.Sprintf("No portfolio found with %s Id", unknownID)
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

	// throws no error when provided oid is not found
	resBool, err = db.DeleteOperation(pid.Hex(), unknownID)
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
	resBool, err = db.DeleteOperation(pid.Hex(), remOp)
	if err != nil {
		t.Errorf("Fail! Unknown error: '%s'", err)
	}
	if resBool != true {
		t.Errorf("Fail! Expected 'true' got '%d'", resInt)
	} else {
		t.Logf("Success! Expected 'true' got '%d'", resInt)
	}

	// we should get one operation after another one has been deleted
	res, err = db.GetOperations(pid.Hex(), "", "", "", "")
	if err != nil {
		t.Errorf("Fail! Unknown error: '%s'", err)
	}
	if len(res) != 1 {
		t.Errorf("Fail! Expected '1' got '%v'", len(res))
	} else {
		t.Logf("Success! Expected '1' got '%v'", len(res))
	}

	// ensure we successfully removed all operations
	resInt, err = db.DeleteOperations(pid.Hex())
	if err != nil {
		t.Errorf("Fail! Unknown error: '%s'", err)
	}
	if resInt != 1 {
		t.Errorf("Fail! Expected '1' got '%d'", resInt)
	} else {
		t.Logf("Success! Expected '1' got '%d'", resInt)
	}

	// we should get no operations
	res, err = db.GetOperations(pid.Hex(), "", "", "", "")
	if err != nil {
		t.Errorf("Fail! Unknown error: '%s'", err)
	}
	if len(res) != 0 {
		t.Errorf("Fail! Expected '0' got '%v'", len(res))
	} else {
		t.Logf("Success! Expected '0' got '%v'", len(res))
	}

	removeTestPortfolio(pid)
	removeTestOperations(pid)
}

func TestSaveSingleOperation(t *testing.T) {
	pid := addTestPortfolio()
	now := getTime()
	ops := getOperations(now)

	// ensure we can get back inserted operation
	_, err := db.AddOperation(pid.Hex(), ops[0])
	if err != nil {
		t.Errorf("Fail! Unknown error: '%s'", err)
	}

	res, err := db.GetOperations(pid.Hex(), "", "", "", "")
	if err != nil {
		t.Errorf("Fail! Unknown error: '%s'", err)
	}
	if ops[0].DateTime != res[0].DateTime {
		t.Errorf("Fail! Expected '%v' got '%v'", ops[0].DateTime, res[0].DateTime)
	} else {
		t.Logf("Success! Expected '%v' got '%v'", ops[0].DateTime, res[0].DateTime)
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

func removeTestPortfolios(pid primitive.ObjectID) {
	ctx := db.context()
	filter := bson.M{"uid": pid}
	opts := options.Delete()

	db.portfolios.DeleteMany(ctx, filter, opts)
}

func addTestUser() primitive.ObjectID {
	ctx := db.context()
	opts := options.InsertOne()
	var testItem interface{} = struct {
		TestKey string `bson:"test_key"`
	}{"test_value"}

	r, _ := db.users.InsertOne(ctx, testItem, opts)
	return r.InsertedID.(primitive.ObjectID)
}

func removeTestUser(id primitive.ObjectID) {
	ctx := db.context()
	filter := bson.M{"_id": id}
	opts := options.Delete()

	db.users.DeleteOne(ctx, filter, opts)
}

func removeTestOperations(pid primitive.ObjectID) {
	ctx := db.context()
	filter := bson.M{"pid": pid}
	opts := options.Delete()

	db.operations.DeleteMany(ctx, filter, opts)
}

func getTime() time.Time {
	t, _ := time.Parse(time.RFC3339, "2020-05-13T22:08:41Z")
	return t
}

func getOperations(now time.Time) []models.Operation {
	ops := []models.Operation{}
	op1 := models.Operation{
		Currency:      models.RUB,
		Price:         50.5351,
		Volume:        150,
		ISIN:          "IE00BD3QHZ91",
		Ticker:        "FXUS",
		DateTime:      now,
		OperationType: models.Buy,
	}

	op2 := models.Operation{
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
