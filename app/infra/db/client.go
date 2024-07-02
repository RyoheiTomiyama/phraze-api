package db

import (
	"context"
	"database/sql"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type client struct {
	db *sqlx.DB
}

type IClient interface {
	GetDeck(ctx context.Context, id int64) (*domain.Deck, error)
}

func NewClient(dataSource string) (IClient, error) {
	db, err := open(dataSource)
	if err != nil {
		return nil, err
	}

	return &client{db}, nil
}

func open(dataSource string) (*sqlx.DB, error) {
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
