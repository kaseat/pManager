package sberbank

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/kaseat/pManager/models"
	"github.com/kaseat/pManager/models/currency"
	"github.com/kaseat/pManager/models/operation"
)

func processSecuritiesInfoTable(rt [][]string) map[ticker]securitiesInfo {
	si := map[ticker]securitiesInfo{}
	if len(rt) > 0 {
		for _, s := range rt[1:] {
			isBond := false
			if s[4] == "Облигация" {
				isBond = true
			}
			si[ticker(s[1])] = securitiesInfo{Ticker: s[1], ISIN: s[2], IsBond: isBond}
		}
	}
	return si
}

func processBuybacksTable(table [][]string, si map[ticker]securitiesInfo) []models.Operation {
	result := []models.Operation{}
	if len(table) > 0 {
		for _, col := range table[3:] {
			if col[3] == "Выкуп ЦБ" {
				rawTime := fmt.Sprintf("%sT10:00:00+03:00", col[0])
				time, _ := time.Parse("02.01.2006T15:04:05Z07:00", rawTime)
				buyback := models.Operation{
					Currency:      currency.RUB,
					Price:         0,
					Volume:        int64(parseFloat(col[5])),
					Ticker:        "RUB",
					DateTime:      time,
					OperationType: operation.Sell,
				}
				if val, ok := si[ticker(col[4])]; ok {
					buyback.ISIN = val.ISIN
				}
				result = append(result, buyback)
			}
		}
	}
	return result
}

func processOperationsTable(table [][]string, si map[ticker]securitiesInfo) []models.Operation {
	result := []models.Operation{}
	if len(table) > 0 {
		for _, row := range table[1:] {
			if len(row) != 16 {
				continue
			}
			if !strings.Contains(row[15], "З") {
				continue
			}
			rawTime := fmt.Sprintf("%sT%s+03:00", row[0], row[2])
			opTime, _ := time.Parse("02.01.2006T15:04:05Z07:00", rawTime)
			var opType operation.Type
			switch o := row[6]; o {
			case "Покупка":
				opType = operation.Buy
			case "Продажа":
				opType = operation.Sell
			default:
				opType = operation.Unknown
			}

			op := models.Operation{
				Currency:      currency.RUB,
				Volume:        int64(parseFloat(row[7])),
				ISIN:          si[ticker(row[4])].ISIN,
				Ticker:        row[4],
				DateTime:      opTime,
				OperationType: opType,
			}

			if si[ticker(row[4])].IsBond {
				op.Price = parseFloat(row[9]) / parseFloat(row[7])

				interest := models.Operation{
					Currency: currency.RUB,
					Price:    parseFloat(row[10]),
					Volume:   1,
					Ticker:   "RUB",
					DateTime: opTime,
				}
				if opType == operation.Buy {
					interest.OperationType = operation.AccInterestBuy
				}
				if opType == operation.Sell {
					interest.OperationType = operation.AccInterestSell
				}
				result = append(result, interest)
			} else {
				op.Price = parseFloat(row[8])
			}

			brokerageFee := models.Operation{
				Currency:      currency.RUB,
				Price:         parseFloat(row[11]),
				Volume:        1,
				Ticker:        "RUB",
				DateTime:      opTime,
				OperationType: operation.BrokerageFee,
			}

			exchangeFee := models.Operation{
				Currency:      currency.RUB,
				Price:         parseFloat(row[12]),
				Volume:        1,
				Ticker:        "RUB",
				DateTime:      opTime,
				OperationType: operation.ExchangeFee,
			}
			result = append(result, op, brokerageFee, exchangeFee)
		}
	}
	return result
}

func processCashFlowTable(rt [][]string) []models.Operation {
	result := []models.Operation{}
	for _, row := range rt[1:] {
		if len(row) != 6 {
			continue
		}
		ind := strings.Index(row[2], "ISIN")
		if ind >= 0 {
			fmt.Println(row[2][ind+5 : ind+17])

			rawTime := fmt.Sprintf("%sT10:00:00+03:00", row[0])
			time, _ := time.Parse("02.01.2006T15:04:05Z07:00", rawTime)
			isin := row[2][ind+5 : ind+17]

			buyback := models.Operation{
				Currency:      currency.Type(row[3]),
				Price:         parseFloat(row[4]),
				Volume:        1,
				ISIN:          isin,
				DateTime:      time,
				OperationType: operation.Buyback,
			}
			result = append(result, buyback)
		}
		if row[2] == "Зачисление д/с" {
			rawTime := fmt.Sprintf("%sT10:00:00+03:00", row[0])
			time, _ := time.Parse("02.01.2006T15:04:05Z07:00", rawTime)

			payIn := models.Operation{
				Currency:      currency.Type(row[3]),
				Price:         1,
				Volume:        int64(parseFloat(row[4])),
				Ticker:        "RUB",
				DateTime:      time,
				OperationType: operation.PayIn,
			}
			result = append(result, payIn)
		}
	}
	return result
}

func parseFloat(str string) float64 {
	var sb strings.Builder
	sb.Grow(len(str))
	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			sb.WriteRune(ch)
		}
	}
	r, _ := strconv.ParseFloat(sb.String(), 64)
	return r
}
