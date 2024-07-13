package db

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/infra/db/model"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
)

func (c *client) CreateCard(ctx context.Context, card *domain.Card) (*domain.Card, error) {
	e := c.execerFrom(ctx)

	query := `INSERT INTO cards (deck_id, question, answer) VALUES (:deck_id, :question, :answer) RETURNING *`
	d := model.Card{DeckID: card.DeckID, Question: card.Question, Answer: card.Answer}

	query, args, err := e.BindNamed(query, d)
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	if err = e.QueryRowxContext(ctx, query, args...).StructScan(&d); err != nil {
		return nil, errutil.Wrap(err)
	}

	return d.ToDomain(), nil
}
