package db_test

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/RyoheiTomiyama/phraze-api/util/env"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

const TestDBName = "phraze_test"

func GetDB() (*sqlx.DB, error) {
	config, err := env.New()
	if err != nil {
		return nil, err
	}

	dataSource := fmt.Sprintf(
		"host=%s port=%s dbname=%s sslmode=disable user=%s password=%s",
		config.DB.HOST, config.DB.PORT, TestDBName, config.DB.USER, config.DB.PASSWORD,
	)

	sqlDB, err := sql.Open("pgx", dataSource)
	db := sqlx.NewDb(sqlDB, "pgx")
	if err != nil {
		return nil, err
	}

	return db, nil
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

	if _, err := db.Exec(`DROP DATABASE IF EXISTS "phraze_test"`); err != nil {
		return fmt.Errorf("could not create test database: %w", err)
	}

	if _, err := db.Exec(`CREATE DATABASE "phraze_test"`); err != nil {
		return fmt.Errorf("could not create test database: %w", err)
	}

	return nil
}

func migration() error {
	migrateFilePath := "../../../atlas/schema.sql"

	config, err := env.New()
	if err != nil {
		return err
	}

	dataSource := fmt.Sprintf(
		"host=%s port=%s dbname=%s sslmode=disable user=%s password=%s",
		config.DB.HOST, config.DB.PORT, TestDBName, config.DB.USER, config.DB.PASSWORD,
	)

	sqlDB, err := sql.Open("pgx", dataSource)
	db := sqlx.NewDb(sqlDB, "pgx")
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed ping: %w", err)
	}

	content, err := os.ReadFile(migrateFilePath)
	if err != nil {
		return fmt.Errorf("could not read schema file: %w", err)
	}

	if _, err := db.Exec(string(content)); err != nil {
		return fmt.Errorf("could not execute SQL %s: %w", string(content), err)
	}

	return nil
}
