// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type Card struct {
	ID int64 `json:"id"`
	// Deck ID
	DeckID int64 `json:"deckId"`
	// 英単語・フレーズ
	Question string `json:"question"`
	// 解答・意味のマークダウン文字列
	Answer    string    `json:"answer"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CardsInput struct {
	Where  *CardsWhere `json:"where"`
	Limit  *int        `json:"limit,omitempty"`
	Offset *int        `json:"offset,omitempty"`
}

type CardsOutput struct {
	Cards    []*Card   `json:"cards"`
	PageInfo *PageInfo `json:"pageInfo"`
}

type CardsWhere struct {
	DeckID int64 `json:"deckId"`
}

type CreateCardInput struct {
	DeckID   int64   `json:"deckId"`
	Question string  `json:"question"`
	Answer   *string `json:"answer,omitempty"`
}

type CreateCardOutput struct {
	Card *Card `json:"card"`
}

type CreateDeckInput struct {
	// Deck名
	Name string `json:"name"`
}

type CreateDeckOutput struct {
	Deck *Deck `json:"deck"`
}

type Deck struct {
	ID        int64     `json:"id"`
	UserID    string    `json:"userId"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeckInfo  *DeckInfo `json:"deckInfo"`
}

type DeckInfo struct {
	TotalCardCount   int        `json:"totalCardCount"`
	PendingCardCount int        `json:"pendingCardCount"`
	LearnedCardCount int        `json:"learnedCardCount"`
	ScheduleAt       *time.Time `json:"scheduleAt,omitempty"`
}

type DecksOutput struct {
	Decks []*Deck `json:"decks"`
}

type Health struct {
	Healthy *bool `json:"healthy,omitempty"`
}

type Mutation struct {
}

type PageInfo struct {
	TotalCount int `json:"totalCount"`
}

type PendingCardsInput struct {
	Where  *CardsWhere `json:"where"`
	Limit  *int        `json:"limit,omitempty"`
	Offset *int        `json:"offset,omitempty"`
}

type PendingCardsOutput struct {
	Cards []*Card `json:"cards"`
}

type Query struct {
}

type ReviewCardInput struct {
	CardID int64 `json:"cardId"`
	Grade  int   `json:"grade"`
}

type ReviewCardOutput struct {
	CardID int64 `json:"cardId"`
}

type UpdateCardInput struct {
	ID       int64   `json:"id"`
	Question *string `json:"question,omitempty"`
	Answer   *string `json:"answer,omitempty"`
}

type UpdateCardOutput struct {
	Card *Card `json:"card"`
}

type Role string

const (
	RoleUser Role = "USER"
)

var AllRole = []Role{
	RoleUser,
}

func (e Role) IsValid() bool {
	switch e {
	case RoleUser:
		return true
	}
	return false
}

func (e Role) String() string {
	return string(e)
}

func (e *Role) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Role(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Role", str)
	}
	return nil
}

func (e Role) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
