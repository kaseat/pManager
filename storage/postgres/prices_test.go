package postgres

import (
	"testing"
	"time"

	"github.com/kaseat/pManager/models"
)

func TestPriceStorage(t *testing.T) {
	tp := getTestPrices()
	err := db.AddPrices(tp)
	if err != nil {
		t.Errorf("Fail! Unexpected error while adding prices %v", err)
	}
	ins, err := db.GetPrices("isin", "US8552441094", "", "")
	if err != nil {
		t.Errorf("Fail! Unexpected error while getting prices %v", err)
	}
	if len(ins) == 4 {
		if ins[0] == tp[0] {
			t.Logf("Success! Expected %v, got %v", tp[0], ins[0])
		} else {
			t.Errorf("Fail! Saved and fetched prices not match! Expected %v, got %v", tp[0], ins[0])
		}
	} else {
		t.Errorf("Fail! Expected %v results, got %v", 1, len(ins))
	}
	dl, _ := db.DeletePrices("isin", "US8552441094")
	if dl != 4 {
		t.Errorf("Fail! Expected %v element to be deleted, got %v", 1, dl)
	}

	ins, _ = db.GetPrices("isin", "US8552441094", "", "")
	if len(ins) == 0 {
		t.Logf("Success! Expected %v, got %v", 0, len(ins))
	} else {
		t.Errorf("Fail! Expected %v pricec after delete, got %v", 0, len(ins))
	}

	db.AddPrices(tp)
	ins, _ = db.GetPricesByIsin("US8552441094", "2019-08-22T07:00:00Z", "2019-08-23T07:00:00Z")
	if len(ins) == 2 {
		t.Logf("Success! Expected %v, got %v", 2, len(ins))
	} else {
		t.Errorf("Fail! Expected %v pricec after delete, got %v", 2, len(ins))
	}

	db.DeleteAllPrices()
	ins, _ = db.GetPricesByIsin("US8552441094", "", "")
	if len(ins) == 0 {
		t.Logf("Success! Expected %v, got %v", 0, len(ins))
	} else {
		t.Errorf("Fail! Expected %v pricec after delete, got %v", 0, len(ins))
	}

	db.DeleteAllInstruments()
}

func getTestPrices() []models.Price {
	db.AddInstruments(getTestInstruments())
	rawInstruments, _ := db.GetAllInstruments()
	instrumentMap := make(map[string]models.Instrument, 2)
	for _, instrument := range rawInstruments {
		instrumentMap[instrument.ISIN] = instrument
	}
	testPrices := getTestPricesRaw()
	for i, testPrice := range testPrices {
		testPrices[i].SecID = instrumentMap[testPrice.ISIN].SecID
	}
	return testPrices
}

func getTestPricesRaw() []models.Price {
	t1, _ := time.Parse("2006-01-02T15:04:05Z", "2019-08-21T00:00:00Z")
	t2, _ := time.Parse("2006-01-02T15:04:05Z", "2019-08-22T00:00:00Z")
	t3, _ := time.Parse("2006-01-02T15:04:05Z", "2019-08-23T00:00:00Z")
	t4, _ := time.Parse("2006-01-02T15:04:05Z", "2019-08-26T00:00:00Z")
	return []models.Price{
		{
			ISIN:   "US8552441094",
			Date:   t1,
			Price:  194.2,
			Volume: 223298,
		},
		{
			ISIN:   "US8552441094",
			Date:   t2,
			Price:  195.73,
			Volume: 223298,
		},
		{
			ISIN:   "US8552441094",
			Date:   t3,
			Price:  196.34,
			Volume: 167941,
		},
		{
			ISIN:   "US8552441094",
			Date:   t4,
			Price:  190,
			Volume: 371916,
		},
	}
}
