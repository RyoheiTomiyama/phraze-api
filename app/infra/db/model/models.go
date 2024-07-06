package model

import (
	"database/sql"
	"time"

	"github.com/RyoheiTomiyama/phraze-api/domain"
)

type Card struct {
	ID        int64
	DeckID    int64
	Question  string
	Answer    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Deck struct {
	ID        int64     `db:"id"`
	UserID    string    `db:"user_id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type User struct {
	ID        int64
	Name      string
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
}

func (d *Deck) ToDomain() *domain.Deck {
	return &domain.Deck{
		ID:        d.ID,
		UserID:    d.UserID,
		Name:      d.Name,
		CreateAt:  d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}
