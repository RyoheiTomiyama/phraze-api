package model

import (
	"database/sql"
	"time"
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
