package mongo

import (
	"testing"
	"time"

	"github.com/kaseat/pManager/models"
	"github.com/kaseat/pManager/models/currency"
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

	dl, _ := db.DeleteInstruments("isin", "RU000A1013Y3")
	if dl != 1 {
		t.Errorf("Fail! Expected %v element to be deleted, got %v", 1, dl)
	}

	db.AddInstruments(ti)

	ins, _ = db.GetAllInstruments()
	if len(ins) != 3 {
		t.Errorf("Fail! Expected %v element to be fetced, got %v", 3, len(ins))
	}

	dl, _ = db.DeleteAllInstruments()
	if dl != 3 {
		t.Errorf("Fail! Expected %v element to be deleted, got %v", 3, dl)
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
			Currency: currency.RUB,
			Type:     instrument.Bond,
		},
		{
			ISIN:     "US8552441094",
			FIGI:     "BBG000CTQBF3",
			Ticker:   "SBUX",
			Name:     "Starbucks Corporation",
			Currency: currency.USD,
			Type:     instrument.Stock,
		},
	}
}
