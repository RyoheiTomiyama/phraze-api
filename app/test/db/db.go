package db_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"ariga.io/atlas-go-sdk/atlasexec"
	"github.com/DATA-DOG/go-txdb"
	"github.com/RyoheiTomiyama/phraze-api/util/env"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

const TestDBName = "phraze_test"

func getTestDataSource() (string, error) {
	config, err := env.New()
	if err != nil {
		return "", err
	}

	dataSource := fmt.Sprintf(
		"host=%s port=%s dbname=%s sslmode=disable user=%s password=%s",
		config.DB.HOST, config.DB.PORT, TestDBName, config.DB.USER, config.DB.PASSWORD,
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

// main_test.goから1度だけ叩く
func SetupDB() (*sqlx.DB, error) {
	fmt.Print("initialize DB...")
	if err := initializeDB(); err != nil {
		return nil, err
	}

	fmt.Print("migratie DB...")

	if err := migration(); err != nil {
		return nil, err
	}

	return nil, nil
}

func initializeDB() error {
	config, err := env.New()
	if err != nil {
		return err
	}

	dataSource := fmt.Sprintf(
		"host=%s port=%s dbname=%s sslmode=disable user=%s password=%s",
		config.DB.HOST, config.DB.PORT, config.DB.DB_NAME, config.DB.USER, config.DB.PASSWORD,
	)

	sqlDB, err := sql.Open("pgx", dataSource)
	db := sqlx.NewDb(sqlDB, "pgx")
	if err != nil {
		return err
	}
	defer db.Close()

	var exist bool
	if err := db.QueryRow(fmt.Sprintf(
		"SELECT EXISTS(SELECT datname FROM pg_database WHERE datname='%s')",
		TestDBName,
	)).Scan(&exist); err != nil {
		return fmt.Errorf("failed database exist check: %w", err)
	}

	if exist {
		return nil
	}

	if _, err := db.Exec(fmt.Sprintf("CREATE DATABASE %s", TestDBName)); err != nil {
		return fmt.Errorf("could not create test database: %w", err)
	}

	return nil
}

func migration() error {
	config, err := env.New()
	if err != nil {
		return err
	}

	dataSource := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.DB.USER, config.DB.PASSWORD, config.DB.HOST, config.DB.PORT, TestDBName,
	)

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	client, err := atlasexec.NewClient(filepath.Join(dir, "../../.."), "atlas")
	if err != nil {
		return err
	}

	_, err = client.SchemaApply(context.Background(), &atlasexec.SchemaApplyParams{
		DevURL: "docker://postgres",
		URL:    dataSource,
		To:     "file://atlas/schema.sql",
	})
	if err != nil {
		return err
	}

	return nil
}
