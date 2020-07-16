package sync

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/kaseat/pManager/gmail"
	"github.com/kaseat/pManager/models"
	"github.com/kaseat/pManager/models/currency"
	"github.com/kaseat/pManager/models/operation"
	"github.com/kaseat/pManager/storage"
)

type pair struct {
	begin, end int
}

type operationInfo struct {
	Currency      string
	Price         float64
	Volume        int64
	ISIN          string
	Ticker        string
	OperationTime time.Time
	OperationType string
}

type securitiesInfo struct {
	Ticker string
	ISIN   string
	IsBond bool
}

// Sberbank sync sber
func Sberbank(login, pid, from, to string) error {
	cl := gmail.GetClient()
	srv, err := cl.GetServiceForUser(login)
	if err != nil {
		return err
	}

	s := storage.GetStorage()
	t, err := s.GetUserLastUpdateTime(login, "sberbank")
	if err != nil {
		return err
	}

	query := "from:broker_rep@sberbank.ru subject:report filename:html"
	if !t.IsZero() {
		query = fmt.Sprintf("%s after:%s", query, t.Format("2006/01/02"))
	}
	if from != "" {
		query = fmt.Sprintf("%s after:%s", query, from)
	}
	if to != "" {
		query = fmt.Sprintf("%s before:%s", query, to)
	}

	r, err := srv.Users.Messages.List("me").Q(query).Do()
	if err != nil {
		return err
	}

	parsedDates := make(map[string]bool)
	operations := make(map[models.Operation]bool)

	for _, m := range r.Messages {
		msg, err := srv.Users.Messages.Get("me", m.Id).Do()

		if err != nil {
			return err
		}
		attachmentID := ""
		for _, p := range msg.Payload.Parts {
			if strings.Contains(p.Filename, ".html") {
				attachmentID = p.Body.AttachmentId
			}
		}

		att, err := srv.Users.Messages.Attachments.Get("me", m.Id, attachmentID).Do()
		if err != nil {
			return err
		}

		b, err := base64.URLEncoding.DecodeString(att.Data)
		if err != nil {
			return err
		}

		ops, err := getOperations(pid, b, parsedDates)
		if err != nil {
			return err
		}

		for _, it := range ops {
			if !operations[it] {
				operations[it] = true
			}
		}

		fmt.Println("Parse message on", time.Unix(msg.InternalDate/1000, 0), "ok!")
	}

	if len(operations) != 0 {
		ops := make([]models.Operation, 0)
		for it := range operations {
			ops = append(ops, it)
		}

		sort.Sort(models.OperationSorter(ops))
		_, err = s.AddOperations(pid, ops)
		if err != nil {
			return err
		}
		fmt.Println("save opertions for", login, "to storge OK")
	}

	err = s.AddUserLastUpdateTime(login, "sberbank", time.Now())
	if err != nil {
		return err
	}
	return nil
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

func findBuyback(table [][]string, si map[string]securitiesInfo) []models.Operation {
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
					ISIN:          si[col[4]].ISIN,
					Ticker:        "RUB",
					DateTime:      time,
					OperationType: operation.Sell,
				}
				result = append(result, buyback)
			}
		}
	}
	return result
}

func findPayIn(rt [][]string, pid string) []models.Operation {
	operations := []models.Operation{}
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
				Ticker:        "RUB",
				DateTime:      time,
				OperationType: operation.Buyback,
			}

			s := storage.GetStorage()
			ops, _ := s.GetOperations(pid, "isin", isin, "", "")

			if len(ops) > 0 {
				buyback.Ticker = ops[0].Ticker
			}

			operations = append(operations, buyback)
		}
		if row[2] == "Зачисление д/с" {
			rawTime := fmt.Sprintf("%sT10:00:00+03:00", row[0])
			time, _ := time.Parse("02.01.2006T15:04:05Z07:00", rawTime)

			payIn := models.Operation{
				Currency:      currency.Type(row[3]),
				Price:         1,
				Volume:        int64(parseFloat(row[4])),
				ISIN:          "BBG0013HGFT4",
				Ticker:        "RUB",
				DateTime:      time,
				OperationType: operation.PayIn,
			}
			operations = append(operations, payIn)
		}
	}
	return operations
}

func getOperations(pid string, b []byte, parsedDates map[string]bool) ([]models.Operation, error) {
	reader := bytes.NewReader(b)
	tables, err := fetchTables(reader, parsedDates)

	if err != nil {
		return nil, err
	}
	if len(tables) == 0 {
		return []models.Operation{}, nil
	}

	operations := parseTable(tables[0])
	securitiesInfo := parseTable(tables[1])
	movements := parseTable(tables[2])
	buybackTable := parseTable(tables[3])
	si := getSecuritiesInfo(securitiesInfo)
	buybacks := findBuyback(buybackTable, si)

	result := []models.Operation{}
	if len(operations) > 0 {
		for _, row := range operations[1:] {
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
				ISIN:          si[row[4]].ISIN,
				Ticker:        row[4],
				DateTime:      opTime,
				OperationType: opType,
			}

			if si[row[4]].IsBond {
				op.Price = parseFloat(row[9]) / parseFloat(row[7])

				interest := models.Operation{
					Currency: currency.RUB,
					Price:    parseFloat(row[10]),
					Volume:   1,
					ISIN:     "BBG0013HGFT4",
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

			brokerFee := models.Operation{
				Currency:      currency.RUB,
				Price:         parseFloat(row[11]),
				Volume:        1,
				ISIN:          "BBG0013HGFT4",
				Ticker:        "RUB",
				DateTime:      opTime,
				OperationType: operation.BrokerageFee,
			}

			exchangeFee := models.Operation{
				Currency:      currency.RUB,
				Price:         parseFloat(row[12]),
				Volume:        1,
				ISIN:          "BBG0013HGFT4",
				Ticker:        "RUB",
				DateTime:      opTime,
				OperationType: operation.ExchangeFee,
			}
			result = append(result, op, brokerFee, exchangeFee)
		}
	}

	if len(movements) > 0 {
		result = append(result, findPayIn(movements, pid)...)
	}

	if len(buybacks) > 0 {
		result = append(result, buybacks...)
	}

	return result, nil
}

func getSecuritiesInfo(rt [][]string) map[string]securitiesInfo {
	si := map[string]securitiesInfo{}
	if len(rt) > 0 {
		for _, s := range rt[1:] {
			isBond := false
			if s[4] == "Облигация" {
				isBond = true
			}
			si[s[1]] = securitiesInfo{Ticker: s[1], ISIN: s[2], IsBond: isBond}
		}
	}
	return si
}

func fetchTables(r io.Reader, parsedDates map[string]bool) ([]string, error) {
	result := make([]string, 4)
	var operationsTable strings.Builder
	var securitiesInfoTable strings.Builder
	var movementsTable strings.Builder
	var buybackTable strings.Builder
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		// discard monthly reports and duplicates
		if scanner.Text() == "<br>Отчет брокера</br>" {
			scanner.Scan()
			re := regexp.MustCompile(`\d{2}.\d{2}.\d{4}`)
			match := re.FindAllString(scanner.Text(), -1)
			if match[0] != match[1] {
				return []string{}, nil
			}
			if parsedDates[match[0]] {
				return []string{}, nil
			}
			parsedDates[match[0]] = true
		}
		if scanner.Text() == "<br>Сделки купли/продажи ценных бумаг</br>" {
			scanner.Scan()
			for scanner.Scan() {
				operationsTable.Write(scanner.Bytes())
				if scanner.Text() == "</table>" {
					break
				}
			}
		}
		if scanner.Text() == "<br>Справочник Ценных Бумаг**</br>" {
			scanner.Scan()
			for scanner.Scan() {
				securitiesInfoTable.Write(scanner.Bytes())
				if scanner.Text() == "</table>" {
					break
				}
			}
		}
		if scanner.Text() == "<br>Движение денежных средств за период</br>" {
			scanner.Scan()
			for scanner.Scan() {
				movementsTable.Write(scanner.Bytes())
				if scanner.Text() == "</table>" {
					break
				}
			}
		}
		if scanner.Text() == "<br>Движение ЦБ, не связанное с исполнением сделок</br>" {
			scanner.Scan()
			for scanner.Scan() {
				buybackTable.Write(scanner.Bytes())
				if scanner.Text() == "</table>" {
					break
				}
			}
		}

	}
	result[0] = operationsTable.String()
	result[1] = securitiesInfoTable.String()
	result[2] = movementsTable.String()
	result[3] = buybackTable.String()
	return result, scanner.Err()
}

func parseTable(rawTable string) [][]string {
	slice := [][]string{}
	dict := map[string]pair{}
	begin, end, n := 0, 0, 0

	for i, char := range rawTable {
		if char == 60 { // rune "<"
			begin = i + 1
			n = 0
		}

		if char == 32 { // rune " "
			n++
			if n == 1 {
				end = i
			}
		}

		if char == 62 { // rune ">"
			if n == 0 {
				end = i
			}
			rawTag := rawTable[begin:end]
			if []rune(rawTag)[0] == 47 { // closing tag
				tag := rawTag[1:] // remove "/" from tag name
				bounds := dict[tag]
				bounds.end = begin - 1
				dict[tag] = bounds

				if tag == "td" { // process delimiter logic
					lastElement := len(slice) - 1
					arr := slice[lastElement]
					slice[lastElement] = append(arr, rawTable[bounds.begin:bounds.end])
				}
			} else { // opening tag
				bounds := dict[rawTag]
				bounds.begin = i + 1
				dict[rawTag] = bounds

				if rawTag == "tr" { // process delimiter logic
					slice = append(slice, []string{})
				}
			}
		}
	}
	return slice
}
