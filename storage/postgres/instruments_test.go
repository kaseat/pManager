package postgres

import (
	"testing"
	"time"

	"github.com/kaseat/pManager/models"
	"github.com/kaseat/pManager/models/currency"
	"github.com/kaseat/pManager/models/exchange"
	"github.com/kaseat/pManager/models/instrument"
)

func TestInstrumentStorage(t *testing.T) {
	ti := getTestInstruments()
	db.AddInstruments(ti)
	ins, _ := db.GetInstruments("isin", "RU000A1013Y3")
	if len(ins) == 1 {
		if ins[0] == ti[0] {
			t.Logf("Success! Expected %v, got %v", ti[0], ins[0])
		} else {
			t.Errorf("Fail! Saved and fetched instruments not match! Expected %v, got %v", ti[0], ins[0])
		}
	} else {
		t.Errorf("Fail! Expected %v results, got %v", 1, len(ins))
	}

	testDate := time.Date(2020, 3, 11, 0, 0, 0, 0, time.UTC)
	db.SetInstrumentPriceUptdTime("RU000A1013Y3", testDate)

	ins, _ = db.GetInstruments("isin", "RU000A1013Y3")
	if len(ins) == 1 {
		if ins[0].PriceUptdTime == testDate {
			t.Logf("Success! Expected %v, got %v", testDate, ins[0].PriceUptdTime)
		} else {
			t.Errorf("Fail! Saved and fetched instruments not match! Expected %v, got %v", testDate, ins[0].PriceUptdTime)
		}
	} else {
		t.Errorf("Fail! Expected %v results, got %v", 1, len(ins))
	}

	db.ClearInstrumentPriceUptdTime("RU000A1013Y3")

	ins, _ = db.GetInstruments("isin", "RU000A1013Y3")
	if len(ins) == 1 {
		if ins[0].PriceUptdTime.IsZero() {
			t.Log("Success! Expected zero time")
		} else {
			t.Errorf("Fail! Expected zero time, got %v", ins[0].PriceUptdTime)
		}
	} else {
		t.Errorf("Fail! Expected %v results, got %v", 1, len(ins))
	}

	db.SetInstrumentPriceUptdTime("RU000A1013Y3", testDate)
	db.SetInstrumentPriceUptdTime("US8552441094", testDate)

	ins, _ = db.GetInstruments("isin", "RU000A1013Y3")

	for _, item := range ins {
		if item.PriceUptdTime == testDate {
			t.Logf("Success! Expected %v, got %v", testDate, item.PriceUptdTime)
		} else {
			t.Errorf("Fail! Saved and fetched instruments not match! Expected %v, got %v", testDate, item.PriceUptdTime)
		}
	}

	db.ClearAllInstrumentPriceUptdTime()
	ins, _ = db.GetInstruments("isin", "RU000A1013Y3")

	for _, item := range ins {
		if item.PriceUptdTime.IsZero() {
			t.Log("Success! Expected zero time")
		} else {
			t.Errorf("Fail! Expected zero time, got %v", ins[0].PriceUptdTime)
		}
	}

	hasDeleted, _ := db.ClearAllInstrumentPriceUptdTime()
	if hasDeleted {
		t.Error("Fail! Expected no elements to be cleares, got some")
	}

	hasDeleted, _ = db.ClearInstrumentPriceUptdTime("RU000A1013Y3")
	if hasDeleted {
		t.Error("Fail! Expected no elements to be cleares, got some")
	}

	dl, _ := db.DeleteInstruments("isin", "RU000A1013Y3")
	if dl != 1 {
		t.Errorf("Fail! Expected %v element to be deleted, got %v", 1, dl)
	}

	dl, _ = db.DeleteAllInstruments()
	if dl != 1 {
		t.Errorf("Fail! Expected %v element to be deleted, got %v", 1, dl)
	}

	ins, _ = db.GetAllInstruments()
	if len(ins) != 0 {
		t.Errorf("Fail! Expected %v element to be fetced, got %v", 0, len(ins))
	}
}

func getTestInstruments() []models.Instrument {
	return []models.Instrument{
		{
			ISIN:     "RU000A1013Y3",
			FIGI:     "BBG00R05JT04",
			Ticker:   "RU000A1013Y3",
			Name:     "Черкизово выпуск 2",
			Exchange: exchange.MOEX,
			Currency: currency.RUB,
			Type:     instrument.Bond,
		},
		{
			ISIN:     "US8552441094",
			FIGI:     "BBG000CTQBF3",
			Ticker:   "SBUX",
			Name:     "Starbucks Corporation",
			Exchange: exchange.SPBEX,
			Currency: currency.USD,
			Type:     instrument.Stock,
		},
	}
}
