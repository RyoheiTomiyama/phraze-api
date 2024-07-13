package domain

import "time"

type Deck struct {
	ID        int64
	UserID    string
	Name      string
	CreateAt  time.Time
	UpdatedAt time.Time
}

type Card struct {
	ID        int64
	DeckID    int64
	Question  string
	Answer    string
	CreateAt  time.Time
	UpdatedAt time.Time
}
