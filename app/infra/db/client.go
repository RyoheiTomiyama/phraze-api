package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type client struct {
	db *sqlx.DB
}

type IClient interface {
	CreateCard(ctx context.Context, card *domain.Card) (*domain.Card, error)
	GetCard(ctx context.Context, id int64) (*domain.Card, error)
	GetCards(ctx context.Context, where *domain.CardsWhere, limit, offset int) ([]*domain.Card, error)
	UpdateCardByID(ctx context.Context, id int64, input *domain.UpdateCardInput) (*domain.Card, error)
	CountCards(ctx context.Context, where *domain.CardsWhere) (*int, error)
	CreateDeck(ctx context.Context, deck *domain.Deck) (*domain.Deck, error)
	GetDeck(ctx context.Context, id int64) (*domain.Deck, error)
	GetDecks(ctx context.Context, userID string) ([]*domain.Deck, error)
}

type DataSourceOption struct {
	Host     string
	Port     string
	DBName   string
	User     string
	Password string
}

func NewClient(opt DataSourceOption) (IClient, error) {
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

func NewTestClient(t *testing.T, db *sqlx.DB) IClient {
	return &client{db}
}
