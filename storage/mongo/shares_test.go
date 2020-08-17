package mongo

import (
	"testing"
	"time"

	"github.com/kaseat/pManager/models"
	"github.com/kaseat/pManager/models/currency"
	"github.com/kaseat/pManager/models/instrument"
	"github.com/kaseat/pManager/models/operation"
)

func TestPortfolioGetShares(t *testing.T) {
	pid := addTestPortfolio()
	db.AddOperations(pid.Hex(), getOperationsForShares())
	db.AddInstruments(getinstrumentsForShares())
	db.AddPrices(getPricesForShres())

	sh, _ := db.GetShares(pid.Hex(), "2019-01-26T07:00:00Z")
	if len(sh) == 3 {
		t.Logf("Success! Expected %v, got %v", 3, len(sh))

		for _, s := range sh {
			switch s.Ticker {
			case "RUB":
				if s.Price == 72177 {
					t.Logf("Success! Expected %v balance on 2019-01-26, got %v", 72177, s.Price)
				} else {
					t.Errorf("Fail! Expected %v balance on 2019-01-26, got %v", 72177, s.Price)
				}
			case "FXIT":
				if s.Price == 4265 {
					t.Logf("Success! Expected %v FXIT price on 2019-01-26, got %v", 4265, s.Price)
				} else {
					t.Errorf("Fail! Expected %v FXIT price on 2019-01-26, got %v", 4265, s.Price)
				}
			case "FXGD":
				if s.Price == 599.5 {
					t.Logf("Success! Expected %v FXGD price on 2019-01-26, got %v", 599.5, s.Price)
				} else {
					t.Errorf("Fail! Expected %v FXGD price on 2019-01-26, got %v", 599.5, s.Price)
				}
			default:
				t.Errorf("Fail! Expected FXGD, FXIT or RUB, got nothing.")
			}
		}
	} else {
		t.Errorf("Fail! Expected %v securities on 2019-01-26, got %v", 3, len(sh))
	}

	db.DeleteAllInstruments()
	db.DeleteAllPrices()
	removeTestOperations(pid)
	removeTestPortfolio(pid)
}

func getOperationsForShares() []models.Operation {
	base, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2019-01-24T07:00:00Z")
	return []models.Operation{
		{Currency: currency.RUB, Price: 1, Volume: 100000, Ticker: "RUB", ISIN: "RUB", DateTime: base, OperationType: operation.PayIn},
		{Currency: currency.RUB, Price: 4195, Volume: 1, Ticker: "FXIT", ISIN: "IE00BD3QJ757", DateTime: base.AddDate(0, 0, 1), OperationType: operation.Buy},
		{Currency: currency.RUB, Price: 590.7, Volume: 40, Ticker: "FXGD", ISIN: "IE00B8XB7377", DateTime: base.AddDate(0, 0, 1), OperationType: operation.Buy},
		{Currency: currency.RUB, Price: 595, Volume: 20, Ticker: "FXGD", ISIN: "IE00B8XB7377", DateTime: base.AddDate(0, 0, 5), OperationType: operation.Sell},
		{Currency: currency.RUB, Price: 4230, Volume: 1, Ticker: "FXIT", ISIN: "IE00BD3QJ757", DateTime: base.AddDate(0, 0, 5), OperationType: operation.Sell},
		{Currency: currency.RUB, Price: 604.9, Volume: 10, Ticker: "FXGD", ISIN: "IE00B8XB7377", DateTime: base.AddDate(0, 0, 6), OperationType: operation.Sell},
	}
}

func getPricesForShres() []models.Price {
	base, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2019-01-24T07:00:00Z")
	return []models.Price{
		{Price: 4189, Volume: 1569, Date: base, ISIN: "IE00BD3QJ757"},
		{Price: 4265, Volume: 1513, Date: base.AddDate(0, 0, 1), ISIN: "IE00BD3QJ757"},
		{Price: 4228, Volume: 1110, Date: base.AddDate(0, 0, 4), ISIN: "IE00BD3QJ757"},
		{Price: 4202, Volume: 3596, Date: base.AddDate(0, 0, 5), ISIN: "IE00BD3QJ757"},
		{Price: 4235, Volume: 3626, Date: base.AddDate(0, 0, 6), ISIN: "IE00BD3QJ757"},
		{Price: 4275, Volume: 6190, Date: base.AddDate(0, 0, 7), ISIN: "IE00BD3QJ757"},
		{Price: 590.5, Volume: 6521, Date: base, ISIN: "IE00B8XB7377"},
		{Price: 599.5, Volume: 10395, Date: base.AddDate(0, 0, 1), ISIN: "IE00B8XB7377"},
		{Price: 603.5, Volume: 12127, Date: base.AddDate(0, 0, 4), ISIN: "IE00B8XB7377"},
		{Price: 606, Volume: 16404, Date: base.AddDate(0, 0, 5), ISIN: "IE00B8XB7377"},
		{Price: 605.5, Volume: 14652, Date: base.AddDate(0, 0, 6), ISIN: "IE00B8XB7377"},
		{Price: 606, Volume: 14653, Date: base.AddDate(0, 0, 7), ISIN: "IE00B8XB7377"},
	}
}

func getinstrumentsForShares() []models.Instrument {
	return []models.Instrument{
		{FIGI: "BBG005DXDPK9", ISIN: "IE00B8XB7377", Ticker: "FXGD", Name: "FinEx Золото", Currency: currency.RUB, Type: instrument.Etf},
		{FIGI: "BBG005HLTYH9", ISIN: "IE00BD3QJ757", Ticker: "FXIT", Name: "FinEx Акции компаний IT-сектора США", Currency: currency.RUB, Type: instrument.Etf},
	}
}
