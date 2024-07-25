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

type CardReview struct {
	ID         int64
	CardID     int64
	ReviewedAt time.Time
	Grade      int64
}

type CardSchedule struct {
	ID         int64
	CardID     int64
	ScheduleAt time.Time
	Interval   int64
	Efactor    float64
}
