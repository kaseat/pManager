package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Init postgresql module
func (db *Db) Init(config Config) error {
	db.context = context.Background()
	conn, err := pgxpool.Connect(db.context, config.ConnString)
	if err != nil {
		return err
	}
	fmt.Println("Init db ok")
	db.connection = conn
	return nil
}

// IsInitialized checks if db initialized
func (db *Db) IsInitialized() bool {
	if db.context == nil {
		return false
	}
	return true
}
