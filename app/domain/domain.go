package domain

import "time"

type Deck struct {
	ID        int64
	UserID    string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Card struct {
	ID        int64
	DeckID    int64
	Question  string
	Answer    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CardReview struct {
	ID         int64
	CardID     int64
	ReviewedAt time.Time
	Grade      int
}

type CardSchedule struct {
	ID         int64
	CardID     int64
	ScheduleAt time.Time
	Interval   int
	Efactor    float64
}
