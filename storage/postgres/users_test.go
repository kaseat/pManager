package postgres

import (
	"testing"
	"time"

	"golang.org/x/oauth2"
)

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

	hasDeleted, err := db.DeleteUser(login)
	if err != nil {
		t.Errorf("Fail! Could not remove test user. Internal error: %s", err)
	}
	if hasDeleted {
		t.Logf("Success! Expected %v, got %v", true, hasDeleted)
	} else {
		t.Errorf("Fail! Did not remove test user! Expected %v, got %v", true, hasDeleted)
	}
}
