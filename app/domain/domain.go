package domain

import "time"

type Deck struct {
	ID        int64
	UserID    string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DeckInfo struct {
	TotalCardCount   int
	PendingCardCount int
	LearnedCardCount int
	ScheduleAt       *time.Time // 学習開始予定時間（Pendingがある場合はnil）
}

type Card struct {
	ID        int64
	DeckID    int64
	Question  string
	Answer    string
	AIAnswer  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CardReview struct {
	ID         int64
	CardID     int64
	UserID     string
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
