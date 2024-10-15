package db

import (
	"context"
	"time"

	"github.com/RyoheiTomiyama/phraze-api/domain"
)

type IClient interface {
	/* cards */
	CreateCard(ctx context.Context, card *domain.Card) (*domain.Card, error)
	GetCard(ctx context.Context, id int64) (*domain.Card, error)
	GetCards(ctx context.Context, where *domain.CardsWhere, limit, offset int) ([]*domain.Card, error)
	GetPendingCards(ctx context.Context, deckID int64, to time.Time, limit, offset int) ([]*domain.Card, error)
	UpdateCardByID(ctx context.Context, id int64, input *domain.UpdateCardInput) (*domain.Card, error)
	CountCards(ctx context.Context, where *domain.CardsWhere) (int, error)
	DeleteCard(ctx context.Context, id int64) (int64, error)

	/* card_reviews */
	CreateCardReview(ctx context.Context, review *domain.CardReview) (*domain.CardReview, error)
	GetLatestCardReview(ctx context.Context, cardID int64) (*domain.CardReview, error)

	/* card_schedules */
	GetCardSchedule(ctx context.Context, cardID int64) (*domain.CardSchedule, error)
	GetLatestCardSchedulesByDeckID(ctx context.Context, deckIDs []int64) (map[int64]*domain.CardSchedule, error)
	GetCardSchedulesByCardID(ctx context.Context, cardIDs []int64) (map[int64]*domain.CardSchedule, error)
	UpsertCardSchedule(ctx context.Context, schedule *domain.CardSchedule) (*domain.CardSchedule, error)

	CreateDeck(ctx context.Context, deck *domain.Deck) (*domain.Deck, error)
	GetDeck(ctx context.Context, id int64) (*domain.Deck, error)
	GetDecks(ctx context.Context, userID string) ([]*domain.Deck, error)
	DeleteDeck(ctx context.Context, id int64) (int64, error)

	GetDeckInfosByDeckID(ctx context.Context, deckIDs []int64) (map[int64]*domain.DeckInfo, error)

	/* permissions */
	GetPermissionsByUserID(ctx context.Context, userID string) ([]*domain.Permission, error)
	HasPermissionByUserID(ctx context.Context, userID string, key domain.PermissionKey) (bool, error)

	Tx(ctx context.Context, fn func(ctx context.Context) error) error
}
