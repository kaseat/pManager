package prices

import (
	"net/http"
	"time"

	"github.com/kaseat/pManager/sync/moex"
	"github.com/kaseat/pManager/sync/spbex"
)

var isSync int32

// Sync starts prices sync
// Sync via MOEX and SPBEX
func Sync() {
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	go moex.Sync("", client)
	go spbex.Sync("", client)
}
