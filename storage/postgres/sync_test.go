package postgres

import (
	"testing"
	"time"

	"github.com/kaseat/pManager/models/provider"
)

func TestLastUpdateTimeStorage(t *testing.T) {
	// arrange
	provider := provider.Sber
	now, _ := time.Parse(time.RFC3339, "2020-05-13T22:08:41Z")
	login := "test_upd_login"
	db.AddUser(login, "a@a.a", "some_hash")

	// ensure we can get back inserted lastUpdateTime
	err := db.AddUserLastUpdateTime(login, provider, now)
	if err != nil {
		t.Errorf("Fail! Could not save '%s' provider. Internal error: %s", provider, err)
	}

	res, err := db.GetUserLastUpdateTime(login, provider)
	if err != nil {
		t.Errorf("Fail! Could not fetch '%s' provider. Internal error: %s", provider, err)
	}
	if res != now {
		t.Errorf("Fail! Saved and fetched time not match! Expected %s, got %s", now, res)
	} else {
		t.Logf("Success! Expected %s, got %s", now, res)
	}

	// check if we actually deleted lastUpdateTime entry
	err = db.DeleteUserLastUpdateTime(login, provider)
	if err != nil {
		t.Errorf("Fail! Could not delete '%s' provider. Internal error: %s", provider, err)
	}

	res, err = db.GetUserLastUpdateTime(login, provider)
	if err != nil {
		t.Errorf("Fail! Could not fetch '%s' provider. Internal error: %s", provider, err)
	}

	if res.IsZero() {
		t.Logf("Success! Expected zero time, got %s", res)
	} else {
		t.Errorf("Fail! Error getting '%s' provider. Expected zero time, got %s", provider, err)
	}

	db.DeleteUser(login)
}
