package model

import (
	"database/sql"
	"time"

	"github.com/RyoheiTomiyama/phraze-api/domain"
)

type Card struct {
	ID        int64     `db:"id"`
	DeckID    int64     `db:"deck_id"`
	Question  string    `db:"question"`
	Answer    string    `db:"answer"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
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

func (d *Card) ToDomain() *domain.Card {
	return &domain.Card{
		ID:        d.ID,
		DeckID:    d.DeckID,
		Question:  d.Question,
		Answer:    d.Answer,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

func (d *Deck) ToDomain() *domain.Deck {
	return &domain.Deck{
		ID:        d.ID,
		UserID:    d.UserID,
		Name:      d.Name,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

type CardSchedule struct {
	ID         int64     `db:"id"`
	CardID     int64     `db:"card_id"`
	ScheduleAt time.Time `db:"schedule_at"`
	Interval   int       `db:"interval"`
	Efactor    float64   `db:"efactor"`
}

func (m *CardSchedule) ToDomain() *domain.CardSchedule {
	return &domain.CardSchedule{
		ID:         m.ID,
		CardID:     m.CardID,
		ScheduleAt: m.ScheduleAt,
		Interval:   m.Interval,
		Efactor:    m.Efactor,
	}
}

type CardReview struct {
	ID         int64     `db:"id"`
	CardID     int64     `db:"card_id"`
	ReviewedAt time.Time `db:"reviewed_at"`
	Grade      int       `db:"grade"`
}

func (m *CardReview) ToDomain() *domain.CardReview {
	return &domain.CardReview{
		ID:         m.ID,
		CardID:     m.CardID,
		ReviewedAt: m.ReviewedAt,
		Grade:      m.Grade,
	}
}
