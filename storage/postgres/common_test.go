package postgres

import (
	"os"
	"testing"
)

var db Db

func TestMain(m *testing.M) {
	db = Db{}
	db.Init(Config{
		ConnString: "host=localhost port=5432 dbname=p_manager_test user=test password=test",
	})

	os.Exit(m.Run())
}
