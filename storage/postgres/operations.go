package postgres

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/kaseat/pManager/models"
)

// AddOperation saves single opertion into a storage
func (db Db) AddOperation(portfolioID string, op models.Operation) (string, error) {
	pid, err := strconv.ParseInt(portfolioID, 10, 32)
	if err != nil {
		return "", errors.New("Invalid portfolio Id format. Expected positive number")
	}

	var id int
	opIds := getOperationTypesByName()
	query := "insert into operations (pid,isin,time,op_id,vol,price) values ($1,$2,$3,$4,$5,$6) returning id;"
	err = db.connection.QueryRow(db.context, query, pid, op.ISIN, op.DateTime.UTC(), opIds[string(op.OperationType)], op.Volume, op.Price).Scan(&id)
	if err != nil {
		pgerr, ok := err.(*pgconn.PgError)
		if !ok {
			return "", err
		}
		if pgerr.Code == "23503" {
			return "", errors.New("could not add operation to unknown portfolio")
		}
		return "", err
	}
	return strconv.Itoa(id), nil
}

// AddOperations saves multiple opertions into a storage
func (db Db) AddOperations(portfolioID string, ops []models.Operation) ([]string, error) {

	pid, err := strconv.ParseInt(portfolioID, 10, 32)
	if err != nil {
		return nil, errors.New("Invalid portfolio Id format. Expected positive number")
	}

	rows := make([][]interface{}, len(ops))
	opIds := getOperationTypesByName()
	for i, op := range ops {
		rows[i] = []interface{}{pid, op.ISIN, op.DateTime.UTC(), opIds[string(op.OperationType)], op.Volume, op.Price}
	}

	colNames := []string{"pid", "isin", "time", "op_id", "vol", "price"}
	_, err = db.connection.CopyFrom(db.context, pgx.Identifier{"operations"}, colNames, pgx.CopyFromRows(rows))
	pgerr, ok := err.(*pgconn.PgError)
	if !ok {
		return nil, err
	}
	if pgerr.Code == "23503" {
		return nil, errors.New("could not add operation to unknown portfolio")
	}
	return nil, err
}

// GetOperations finds operations depending on input prameters
func (db Db) GetOperations(portfolioID string, key string, value string, from string, to string) ([]models.Operation, error) {
	pid, err := strconv.ParseInt(portfolioID, 10, 32)
	if err != nil {
		return nil, errors.New("Invalid portfolio Id format. Expected positive number")
	}
	params := []interface{}{pid}
	n := 1
	query := `select o.id::varchar(20), o.isin, s.figi, s.currency, o.time, t.name, o.vol, o.price from operations o
	inner join securities s on s.isin = o.isin inner join operation_types t on t.id = o.op_id where pid = $1`
	if key != "" && value != "" {
		n++
		query += fmt.Sprintf(" and %s = $%d", key, n)
		params = append(params, value)
	}
	if dtime, err := time.Parse("2006-01-02T15:04:05Z07:00", from); err == nil {
		n++
		query += fmt.Sprintf(" and time >= $%d", n)
		params = append(params, dtime)
	}
	if dtime, err := time.Parse("2006-01-02T15:04:05Z07:00", to); err == nil {
		n++
		query += fmt.Sprintf(" and time <= $%d", n)
		params = append(params, dtime)
	}
	query += ";"

	rows, err := db.connection.Query(db.context, query, params...)
	if err != nil {
		return nil, err
	}
	result := []models.Operation{}
	for rows.Next() {
		op := models.Operation{
			PortfolioID: portfolioID,
		}
		err = rows.Scan(&op.OperationID, &op.ISIN, &op.FIGI, &op.Currency, &op.DateTime, &op.OperationType, &op.Volume, &op.Price)
		if err != nil {
			return nil, err
		}
		result = append(result, op)
	}

	return result, nil
}

// DeleteOperation removes operation by Id
func (db Db) DeleteOperation(portfolioID string, operationID string) (bool, error) {
	pid, err := strconv.ParseInt(portfolioID, 10, 32)
	if err != nil {
		return false, errors.New("Invalid portfolio Id format. Expected positive number")
	}
	id, err := strconv.ParseInt(operationID, 10, 32)
	if err != nil {
		return false, errors.New("Invalid operation Id format. Expected positive number")
	}

	query := "delete from operations where pid = $1 and id = $2;"
	r, err := db.connection.Exec(db.context, query, pid, id)
	if err != nil {
		return false, err
	}
	if r.RowsAffected() == 1 {
		return true, nil
	}
	return false, nil
}

// DeleteOperations removes all operations for provided portfolio Id
func (db Db) DeleteOperations(portfolioID string) (int64, error) {
	pid, err := strconv.ParseInt(portfolioID, 10, 32)
	if err != nil {
		return 0, errors.New("Invalid portfolio Id format. Expected positive number")
	}

	query := "delete from operations where pid = $1;"
	r, err := db.connection.Exec(db.context, query, pid)
	if err != nil {
		return 0, err
	}
	return r.RowsAffected(), nil
}

func getOperationTypesByName() map[string]int {
	return map[string]int{
		"buy":                 1,
		"sell":                2,
		"brokerageFee":        3,
		"exchangeFee":         4,
		"payIn":               5,
		"payOut":              6,
		"coupon":              7,
		"accruedInterestBuy":  8,
		"accruedInterestSell": 9,
		"buyback":             10,
	}
}
