package moex

import (
	"net/http"
	"testing"
	"time"
)

func TestMoexSync(t *testing.T) {

	client := &http.Client{
		Timeout: time.Second * 5,
	}
	Sync("RU000A0JW1K9", client)
	t.Fail()
}
