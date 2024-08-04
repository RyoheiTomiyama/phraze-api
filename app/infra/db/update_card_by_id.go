package db

import (
	"context"
	"strings"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/infra/db/model"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
)

func (c *client) UpdateCardByID(ctx context.Context, id int64, input *domain.UpdateCardInput) (*domain.Card, error) {
	e := c.execerFrom(ctx)

	query := `UPDATE cards `

	var sets []string
	arg := map[string]interface{}{"id": id}
	if input != nil {
		if input.Field.DeckID != nil {
			sets = append(sets, "deck_id = :deck_id")
			arg["deck_id"] = input.Field.DeckID
		}
		if input.Field.Question != nil {
			sets = append(sets, "question = :question")
			arg["question"] = input.Field.Question
		}
		if input.Field.Answer != nil {
			sets = append(sets, "answer = :answer")
			arg["answer"] = input.Field.Answer
		}
		if input.Field.AIAnswer != nil {
			sets = append(sets, "ai_answer = :ai_answer")
			arg["ai_answer"] = input.Field.AIAnswer
		}
	}

	if len(sets) > 0 {
		query = query + " SET " + strings.Join(sets, ", ")
	}
	query = query + " WHERE id = :id RETURNING *"

	query, args, err := e.BindNamed(query, arg)
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	var card model.Card
	if err = e.QueryRowxContext(ctx, query, args...).StructScan(&card); err != nil {
		return nil, errutil.Wrap(err)
	}

	return card.ToDomain(), nil
}
