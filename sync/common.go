package sync

import (
	"github.com/kaseat/pManager/storage"
)

var db storage.Db

// Init initializes sync module
func Init(storage storage.Db) {
	db = storage
}
