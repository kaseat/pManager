package utils

import (
	"math"

	"github.com/kaseat/pManager/models"
	"github.com/kaseat/pManager/models/operation"
	"github.com/oleiade/lane"
)

// GetSum returns blance for given operations
func GetSum(operations []models.Operation) float64 {
	sum := int64(0)
	for _, op := range operations {
		amount := int64(math.Round(op.Price*1e6)) * op.Volume
		switch op.OperationType {
		case operation.PayIn, operation.Buyback, operation.Sell, operation.AccInterestSell:
			sum += amount
		default:
			sum -= amount
		}
	}
	return math.Round(float64(sum)/1e4) / 100
}

// GetAverage returns average price of given operations
func GetAverage(ops []models.Operation) float64 {
	d := lane.NewDeque()
	for _, op := range ops {
		if op.OperationType == operation.Buy {
			d.Append(op)
		} else {
			for {
				if d.Empty() {
					break
				}
				o := d.Shift().(models.Operation)
				if o.Volume-op.Volume <= 0 {
					op.Volume -= o.Volume
				} else {
					o.Volume -= op.Volume
					d.Prepend(o)
					break
				}
			}
		}
	}

	cost, vol := float64(0), float64(0)
	for {
		if d.Empty() {
			break
		}
		op := d.Pop().(models.Operation)
		cost += op.Price * float64(op.Volume)
		vol += float64(op.Volume)
	}

	result := float64(0)
	if vol != 0 {
		result = math.Round(cost/vol*1e6) / 1e6
	}
	return result
}
