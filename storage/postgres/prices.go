package postgres

import (
	"fmt"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/kaseat/pManager/models"
)

// AddPrices saves prices series into a storage
func (db Db) AddPrices(prices []models.Price) error {
	colNames := []string{"isin", "date", "vol", "price"}
	rows := make([][]interface{}, len(prices))
	for i, pr := range prices {
		rows[i] = []interface{}{pr.ISIN, pr.Date, pr.Volume, pr.Price}
	}

	_, err := db.connection.CopyFrom(db.context, pgx.Identifier{"prices"}, colNames, pgx.CopyFromRows(rows))
	pgerr, ok := err.(*pgconn.PgError)
	if !ok {
		return err
	}
	if pgerr.Code == "23503" {
		return fmt.Errorf("could not add prices: error in column %s", pgerr.ColumnName)
	}
	return nil
}

// GetPrices finds prices depending on input prameters
func (db Db) GetPrices(key, value, from, to string) ([]models.Price, error) {
	n := 0
	params := []interface{}{}
	query := "select isin,date,vol,price from prices"
	if key != "" && value != "" {
		n++
		if n == 1 {
			query += fmt.Sprintf(" where %s = $%d", key, n)
		} else {
			query += fmt.Sprintf(" and %s = $%d", key, n)
		}
		params = append(params, value)
	}
	if dtime, err := time.Parse("2006-01-02T15:04:05Z07:00", from); err == nil {
		n++
		if n == 1 {
			query += fmt.Sprintf(" where date >= $%d", n)
		} else {
			query += fmt.Sprintf(" and date >= $%d", n)
		}
		params = append(params, dtime)
	}
	if dtime, err := time.Parse("2006-01-02T15:04:05Z07:00", to); err == nil {
		n++
		if n == 1 {
			query += fmt.Sprintf(" where date <= $%d", n)
		} else {
			query += fmt.Sprintf(" and date <= $%d", n)
		}
		params = append(params, dtime)
	}
	query += ";"

	rows, err := db.connection.Query(db.context, query, params...)
	if err != nil {
		return nil, err
	}
	result := []models.Price{}
	for rows.Next() {
		pr := models.Price{}
		err = rows.Scan(&pr.ISIN, &pr.Date, &pr.Volume, &pr.Price)
		if err != nil {
			return nil, err
		}
		result = append(result, pr)
	}

	return result, nil
}

// GetPricesByIsin finds prices for given ISIN and dates
func (db Db) GetPricesByIsin(isin, from, to string) ([]models.Price, error) {
	return db.GetPrices("isin", isin, from, to)
}

// DeletePrices removes prices depending on input prameters
func (db Db) DeletePrices(key string, value string) (int64, error) {
	var r pgconn.CommandTag
	var err error
	if key != "" && value != "" {
		query := fmt.Sprintf("delete from prices where %s = $1;", key)
		r, err = db.connection.Exec(db.context, query, value)
	} else {
		query := "delete from prices;"
		r, err = db.connection.Exec(db.context, query)
	}
	if err != nil {
		return 0, err
	}
	return r.RowsAffected(), nil
}

// DeleteAllPrices removes all prices from storage
func (db Db) DeleteAllPrices() (int64, error) {
	return db.DeletePrices("", "")
}
