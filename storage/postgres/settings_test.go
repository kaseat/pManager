package postgres

import "testing"

func TestTcsTokenStorage(t *testing.T) {
	token := "test_token"
	err := db.AddTcsToken(token)
	if err != nil {
		t.Errorf("Fail! error during save token: %v", err)
	}
	res, _ := db.GetTcsToken()
	if res == token {
		t.Logf("Success! Expected %v, got %v", token, res)
	} else {
		t.Errorf("Fail! Saved and fetched tokens not match! Expected %v, got %v", token, res)
	}

	db.DeleteTcsToken()

	res, _ = db.GetTcsToken()

	if res == "" {
		t.Logf("Success! Expected empty string, got %v", res)
	} else {
		t.Errorf("Fail! Expected empty string, got %v", res)
	}
}
