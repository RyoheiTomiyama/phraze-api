package db_test

import (
	"context"
	"database/sql"
	"fmt"

	"ariga.io/atlas-go-sdk/atlasexec"
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

	if _, err := db.Exec(fmt.Sprintf(`
		SELECT 'CREATE DATABASE %s'
		WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '%s')
	`, TestDBName, TestDBName)); err != nil {
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

	client, err := atlasexec.NewClient("..", "atlas")
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
