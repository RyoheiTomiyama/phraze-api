package db_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/RyoheiTomiyama/phraze-api/util/env"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func getTestDataSource() (string, error) {
	config, err := env.New()
	if err != nil {
		return "", err
	}

	dataSource := fmt.Sprintf(
		"host=%s port=%s dbname=%s sslmode=disable user=%s password=%s",
		config.DB.HOST, config.DB.PORT, fmt.Sprintf("%s_test", config.DB.DB_NAME), config.DB.USER, config.DB.PASSWORD,
	)

	return dataSource, nil
}

func GetDB(t *testing.T) *sqlx.DB {
	t.Helper()

	ds, err := getTestDataSource()
	if err != nil {
		t.Fatal(err)

		return nil
	}

	txDB := sql.OpenDB(txdb.New("pgx", ds))
	db := sqlx.NewDb(txDB, "pgx")

	if err := db.Ping(); err != nil {
		t.Fatal(err)

		return nil
	}

	return db
}
