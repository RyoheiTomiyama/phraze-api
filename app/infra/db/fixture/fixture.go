package fixture

import (
	"github.com/RyoheiTomiyama/phraze-api/infra/db/model"
	"github.com/jmoiron/sqlx"
)

type fixture struct {
	db            *sqlx.DB
	Cards         []*model.Card
	CardReviews   []*model.CardReview
	CardSchedules []*model.CardSchedule
	Decks         []*model.Deck
}

func New(db *sqlx.DB) *fixture {
	return &fixture{db: db}
}
