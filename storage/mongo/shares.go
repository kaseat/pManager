package mongo

import (
	"github.com/kaseat/pManager/models"
	"github.com/kaseat/pManager/models/operation"
)

// GetShares gets shares
func (db Db) GetShares(pid string, onDate string) ([]models.Share, error) {
	ops, err := db.GetOperations(pid, "", "", "", onDate)
	if err != nil {
		return nil, err
	}
	shares := make(map[string]models.Share)
	for _, op := range ops {
		if i, ok := shares[op.ISIN]; ok {
			if op.OperationType == operation.Buy {
				i.Volume += op.Volume
				shares[op.ISIN] = i
			}
			if op.OperationType == operation.Sell || op.OperationType == operation.Buyback {
				i.Volume -= op.Volume
				shares[op.ISIN] = i
			}
		} else {
			shares[op.ISIN] = models.Share{
				ISIN:   op.ISIN,
				Volume: op.Volume,
			}
		}
	}
	result := []models.Share{}
	for _, sh := range shares {
		if sh.ISIN != "" && sh.Volume != 0 {
			result = append(result, sh)
		}
	}
	return result, nil
}
