package postgres

import (
	"fmt"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/kaseat/pManager/models"
	"github.com/kaseat/pManager/models/exchange"
	"github.com/kaseat/pManager/models/instrument"
)

// AddInstruments saves instruments info into a storage
func (db Db) AddInstruments(instr []models.Instrument) error {
	sType := getsecuritiesTypeByName()
	exType := getExchangeIDByName()
	colNames := []string{"isin", "ticker", "figi", "currency", "exchange_id", "asset_type", "title"}
	rows := make([][]interface{}, len(instr))
	for i, ins := range instr {
		rows[i] = []interface{}{ins.ISIN, ins.Ticker, ins.FIGI, ins.Currency, exType[ins.Exchange], sType[ins.Type], ins.Name}
	}

	_, err := db.connection.CopyFrom(db.context, pgx.Identifier{"securities"}, colNames, pgx.CopyFromRows(rows))
	pgerr, ok := err.(*pgconn.PgError)
	if !ok {
		return err
	}
	if pgerr.Code == "23503" {
		return fmt.Errorf("could not add instrument: error in column %s", pgerr.ColumnName)
	}
	return nil
}

// SetInstrumentPriceUptdTime sets time instrument prise was updated
func (db Db) SetInstrumentPriceUptdTime(isin string, updTime time.Time) (bool, error) {
	query := "update securities set price_upd_time = $1 where isin = $2;"
	r, err := db.connection.Exec(db.context, query, updTime, isin)
	if err != nil {
		return false, err
	}
	if r.RowsAffected() == 0 {
		return false, nil
	}
	return true, nil
}

// ClearInstrumentPriceUptdTime clears time instrument prise was updated
func (db Db) ClearInstrumentPriceUptdTime(isin string) (bool, error) {
	query := "update securities set price_upd_time = NULL where isin = $1 and price_upd_time is not null;"
	r, err := db.connection.Exec(db.context, query, isin)
	if err != nil {
		return false, err
	}
	if r.RowsAffected() == 0 {
		return false, nil
	}
	return true, nil
}

// ClearAllInstrumentPriceUptdTime clears time instrument prise was updated (for all)
func (db Db) ClearAllInstrumentPriceUptdTime() (bool, error) {
	query := "update securities set price_upd_time = NULL where price_upd_time is not null;"
	r, err := db.connection.Exec(db.context, query)
	if err != nil {
		return false, err
	}
	if r.RowsAffected() == 0 {
		return false, nil
	}
	return true, nil
}

// GetInstruments finds instruments depending on input prameters
func (db Db) GetInstruments(key string, value string) ([]models.Instrument, error) {
	var rows pgx.Rows
	var err error
	query := `select isin,ticker,figi,currency,code,id_name,s.title,price_upd_time
		from securities s inner join securities_types t on t.id = s.asset_type
		inner join exchange e on e.id = s.exchange_id`
	if key != "" && value != "" {
		query += fmt.Sprintf(" where %s = $1;", key)
		rows, err = db.connection.Query(db.context, query, value)
	} else {
		query += ";"
		rows, err = db.connection.Query(db.context, query)
	}
	if err != nil {
		return nil, err
	}

	result := []models.Instrument{}
	for rows.Next() {
		ins := models.Instrument{}
		var tm *time.Time
		err = rows.Scan(&ins.ISIN, &ins.Ticker, &ins.FIGI, &ins.Currency, &ins.Exchange, &ins.Type, &ins.Name, &tm)
		if err != nil {
			return nil, err
		}
		if tm != nil {
			ins.PriceUptdTime = *tm
		}
		result = append(result, ins)
	}
	return result, nil
}

// GetAllInstruments finds all instruments
func (db Db) GetAllInstruments() ([]models.Instrument, error) {
	query := `select isin,ticker,figi,currency,id_name,s.title,price_upd_time
		from securities s inner join securities_types t on t.id = s.asset_type`
	rows, err := db.connection.Query(db.context, query)
	if err != nil {
		return nil, err
	}

	result := []models.Instrument{}
	for rows.Next() {
		ins := models.Instrument{}
		var tm *time.Time
		err = rows.Scan(&ins.ISIN, &ins.Ticker, &ins.FIGI, &ins.Currency, &ins.Type, &ins.Name, &tm)
		if err != nil {
			return nil, err
		}
		if tm != nil {
			ins.PriceUptdTime = *tm
		}
		result = append(result, ins)
	}
	return result, nil
}

// DeleteInstruments removes instruments depending on input prameters
func (db Db) DeleteInstruments(key string, value string) (int64, error) {
	var r pgconn.CommandTag
	var err error
	if key != "" && value != "" {
		query := fmt.Sprintf("delete from securities where %s = $1;", key)
		r, err = db.connection.Exec(db.context, query, value)
	} else {
		query := "delete from securities;"
		r, err = db.connection.Exec(db.context, query)
	}
	if err != nil {
		return 0, err
	}
	return r.RowsAffected(), nil
}

// DeleteAllInstruments removes all instruments from storage
func (db Db) DeleteAllInstruments() (int64, error) {
	query := "delete from securities;"
	r, err := db.connection.Exec(db.context, query)
	if err != nil {
		return 0, err
	}
	return r.RowsAffected(), nil
}

func getsecuritiesTypeByName() map[instrument.Type]int {
	return map[instrument.Type]int{
		instrument.Stock:       10,
		instrument.Bond:        20,
		instrument.EtfStock:    31,
		instrument.EtfBond:     32,
		instrument.EtfMixed:    34,
		instrument.EtfGold:     35,
		instrument.EtfCurrency: 36,
		instrument.Currency:    60,
	}
}

func getExchangeIDByName() map[exchange.Type]int {
	return map[exchange.Type]int{
		exchange.MOEX:  1,
		exchange.SPBEX: 2,
	}
}
