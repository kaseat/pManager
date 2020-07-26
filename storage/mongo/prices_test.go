package mongo

import (
	"testing"
	"time"

	"github.com/kaseat/pManager/models"
)

func TestPriceStorage(t *testing.T) {
	tp := getTestPrices()
	db.AddPrices(tp)
	ins, _ := db.GetPrices("isin", "IE00B4BNMY34")
	if len(ins) == 4 {
		if ins[0] == tp[0] {
			t.Logf("Success! Expected %v, got %v", tp[0], ins[0])
		} else {
			t.Errorf("Fail! Saved and fetched prices not match! Expected %v, got %v", tp[0], ins[0])
		}
	} else {
		t.Errorf("Fail! Expected %v results, got %v", 1, len(ins))
	}
	dl, _ := db.DeletePrices("isin", "IE00B4BNMY34")
	if dl != 4 {
		t.Errorf("Fail! Expected %v element to be deleted, got %v", 1, dl)
	}

	ins, _ = db.GetPrices("isin", "IE00B4BNMY34")
	if len(ins) == 0 {
		t.Logf("Success! Expected %v, got %v", 0, len(ins))
	} else {
		t.Errorf("Fail! Expected %v pricec after delete, got %v", 0, len(ins))
	}

	db.AddPrices(tp)
	ins, _ = db.GetPricesByIsin("IE00B4BNMY34", "2019-08-22T07:00:00Z", "2019-08-23T07:00:00Z")
	if len(ins) == 2 {
		t.Logf("Success! Expected %v, got %v", 0, len(ins))
	} else {
		t.Errorf("Fail! Expected %v pricec after delete, got %v", 0, len(ins))
	}

	db.DeleteAllPrices()
	ins, _ = db.GetPricesByIsin("IE00B4BNMY34", "", "")
	if len(ins) == 0 {
		t.Logf("Success! Expected %v, got %v", 0, len(ins))
	} else {
		t.Errorf("Fail! Expected %v pricec after delete, got %v", 0, len(ins))
	}
}

func getTestPrices() []models.Price {
	t1, _ := time.Parse("2006-01-02T15:04:05Z", "2019-08-21T07:00:00Z")
	t2, _ := time.Parse("2006-01-02T15:04:05Z", "2019-08-22T07:00:00Z")
	t3, _ := time.Parse("2006-01-02T15:04:05Z", "2019-08-23T07:00:00Z")
	t4, _ := time.Parse("2006-01-02T15:04:05Z", "2019-08-26T07:00:00Z")
	return []models.Price{
		{
			ISIN:   "IE00B4BNMY34",
			Date:   t1,
			Price:  194.2,
			Volume: 223298,
		},
		{
			ISIN:   "IE00B4BNMY34",
			Date:   t2,
			Price:  195.73,
			Volume: 223298,
		},
		{
			ISIN:   "IE00B4BNMY34",
			Date:   t3,
			Price:  196.34,
			Volume: 167941,
		},
		{
			ISIN:   "IE00B4BNMY34",
			Date:   t4,
			Price:  190,
			Volume: 371916,
		},
	}
}
