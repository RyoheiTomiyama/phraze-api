package model

import (
	"database/sql"
)

type Card struct {
	ID        int64
	DeckID    int64
	Question  string
	Answer    string
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
}

type Deck struct {
	ID        int64        `db:"id"`
	UserID    string       `db:"user_id"`
	Name      string       `db:"name"`
	CreatedAt sql.NullTime `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updatedAt"`
}

type User struct {
	ID        int64
	Name      string
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
}
