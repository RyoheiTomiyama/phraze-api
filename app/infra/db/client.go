package db

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/RyoheiTomiyama/phraze-api/domain/infra/db"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type client struct {
	db *sqlx.DB
}

type DataSourceOption struct {
	Host     string
	Port     string
	DBName   string
	User     string
	Password string
}

func NewClient(opt DataSourceOption) (db.IClient, error) {
	db, err := open(opt)
	if err != nil {
		return nil, err
	}

	return &client{db}, nil
}

func open(opt DataSourceOption) (*sqlx.DB, error) {
	dataSource := fmt.Sprintf(
		"host=%s port=%s dbname=%s sslmode=disable user=%s password=%s",
		opt.Host, opt.Port, opt.DBName, opt.User, opt.Password,
	)

	sqlDB, err := sql.Open("pgx", dataSource)
	db := sqlx.NewDb(sqlDB, "pgx")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return db, err
	}

	return db, nil
}

func NewTestClient(t *testing.T, db *sqlx.DB) db.IClient {
	return &client{db}
}
