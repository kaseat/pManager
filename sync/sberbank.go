package sync

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/kaseat/pManager/models"
	"github.com/kaseat/pManager/storage"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
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
func Sberbank(login string, pid string) error {
	srv, err := getGmailService()
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

	r, err := srv.Users.Messages.List("me").Q(query).Do()
	if err != nil {
		return err
	}

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
		reader := bytes.NewReader(b)
		op, sir, mv, _ := fetchTables(reader)
		opInfo := getOperationsInfo(parseTable(op), parseTable(mv), getSecuritiesInfo(parseTable(sir)))

		ops := make([]models.Operation, len(opInfo))
		for i, o := range opInfo {
			t := models.Operation{
				PortfolioID:   pid,
				Currency:      models.Currency(o.Currency),
				Price:         o.Price,
				Volume:        o.Volume,
				ISIN:          o.ISIN,
				Ticker:        o.Ticker,
				DateTime:      o.OperationTime,
				OperationType: models.OperationType(o.OperationType),
			}
			ops[i] = t
		}

		_, err = s.AddOperations(pid, ops)
		if err != nil {
			return err
		}
	}
	return nil
}

func getGmailService() (*gmail.Service, error) {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		return nil, err
	}
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		return nil, err
	}

	f, err := os.Open("token.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)

	client := config.Client(context.Background(), token)

	srv, err := gmail.New(client)
	if err != nil {
		return nil, err
	}
	return srv, nil
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

func findPayIn(rt [][]string) []operationInfo {
	operations := []operationInfo{}
	for _, row := range rt[1:] {
		if len(row) != 6 {
			continue
		}
		if row[2] != "Зачисление д/с" {
			continue
		}
		rawTime := fmt.Sprintf("%sT10:00:00+03:00", row[0])
		time, _ := time.Parse("02.01.2006T15:04:05Z07:00", rawTime)

		payIn := operationInfo{
			Currency:      row[3],
			Price:         1,
			Volume:        int64(parseFloat(row[4])),
			ISIN:          "BBG0013HGFT4",
			Ticker:        "RUB",
			OperationTime: time,
			OperationType: "payIn",
		}
		operations = append(operations, payIn)
	}
	return operations
}

func getOperationsInfo(rt [][]string, mv [][]string, si map[string]securitiesInfo) []operationInfo {
	operations := []operationInfo{}

	for _, row := range rt[1:] {
		if len(row) != 16 {
			continue
		}
		rawTime := fmt.Sprintf("%sT%s+03:00", row[0], row[2])
		time, _ := time.Parse("02.01.2006T15:04:05Z07:00", rawTime)
		var opType string
		switch o := row[6]; o {
		case "Покупка":
			opType = "buy"
		case "Продажа":
			opType = "sell"
		default:
			opType = "unknown"
		}

		op := operationInfo{
			Currency:      row[5],
			Price:         parseFloat(row[8]),
			Volume:        int64(parseFloat(row[7])),
			ISIN:          si[row[4]].ISIN,
			Ticker:        row[4],
			OperationTime: time,
			OperationType: opType,
		}

		brokerFee := operationInfo{
			Currency:      row[5],
			Price:         -parseFloat(row[11]),
			Volume:        1,
			ISIN:          "BBG0013HGFT4",
			Ticker:        "RUB",
			OperationTime: time,
			OperationType: "brokerFee",
		}

		exchangeFee := operationInfo{
			Currency:      row[5],
			Price:         -parseFloat(row[12]),
			Volume:        1,
			ISIN:          "BBG0013HGFT4",
			Ticker:        "RUB",
			OperationTime: time,
			OperationType: "exchangeFee",
		}
		operations = append(operations, op, brokerFee, exchangeFee)
	}
	operations = append(operations, findPayIn(mv)...)

	return operations
}

func getSecuritiesInfo(rt [][]string) map[string]securitiesInfo {
	si := map[string]securitiesInfo{}

	for _, s := range rt[1:] {
		isBond := false
		if s[4] == "Облигация" {
			isBond = true
		}
		si[s[1]] = securitiesInfo{Ticker: s[1], ISIN: s[2], IsBond: isBond}
	}
	return si
}

func fetchTables(r io.Reader) (string, string, string, error) {
	var operationsTable strings.Builder
	var securitiesInfoTable strings.Builder
	var movementsTable strings.Builder
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
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
	}
	return operationsTable.String(), securitiesInfoTable.String(), movementsTable.String(), scanner.Err()
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
