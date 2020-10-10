package spbex

import (
	"net/http"
	"testing"
	"time"
)

func TestSpbexSync(t *testing.T) {
	client := &http.Client{
		Timeout: time.Second * 5,
	}

	Sync("MSFT", client)
	t.Fail()
}
