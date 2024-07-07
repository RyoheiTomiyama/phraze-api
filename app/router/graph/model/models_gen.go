// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"
)

type Card struct {
	ID        int64     `json:"id"`
	DeckID    int64     `json:"deckID"`
	Question  string    `json:"question"`
	Answer    string    `json:"answer"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Deck struct {
	ID        int64     `json:"id"`
	UserID    string    `json:"userId"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Health struct {
	Healthy *bool `json:"healthy,omitempty"`
}

type Mutation struct {
}

type Query struct {
}
