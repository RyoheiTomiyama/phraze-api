package model

import (
	"time"

	"github.com/RyoheiTomiyama/phraze-api/domain"
)

type Card struct {
	ID        int64     `db:"id"`
	DeckID    int64     `db:"deck_id"`
	Question  string    `db:"question"`
	Answer    string    `db:"answer"`
	AIAnswer  string    `db:"ai_answer"`
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

func (d *Card) ToDomain() *domain.Card {
	return &domain.Card{
		ID:        d.ID,
		DeckID:    d.DeckID,
		Question:  d.Question,
		Answer:    d.Answer,
		AIAnswer:  d.AIAnswer,
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

type Permission struct {
	ID   int64  `db:"id"`
	Key  string `db:"key"`
	Name string `db:"name"`
}

type Role struct {
	ID   int64  `db:"id"`
	Key  string `db:"key"`
	Name string `db:"name"`
}

type RolesPermission struct {
	RoleID       int64 `db:"role_id"`
	PermissionID int64 `db:"permission_id"`
}

type User struct {
	ID        string    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type UsersRole struct {
	UserID    string    `db:"user_id"`
	RoleID    int64     `db:"role_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
