package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Config represents database configuration
type Config struct {
	ConnString string `json:"connString"`
	DbName     string `json:"dbName"`
}

// Db represents storage
type Db struct {
	connection *pgxpool.Pool
	context    context.Context
}
