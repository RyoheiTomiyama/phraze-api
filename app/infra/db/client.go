package db

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/domain"
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
	db, err := sqlx.Open("pgx", dataSource)
	if err != nil {
		return db, err
	}

	if err := db.Ping(); err != nil {
		return db, err
	}

	return db, nil
}
