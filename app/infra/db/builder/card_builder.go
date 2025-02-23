package builder

import (
	"context"
	"fmt"
	"strings"

	"github.com/RyoheiTomiyama/phraze-api/domain"
)

type CardsWhere domain.CardsWhere

func (builder *CardsWhere) BuildNamedWhere(ctx context.Context,
	query string, arg map[string]interface{},
) (string, map[string]interface{}) {
	var wheres []string
	if builder != nil {
		if builder.DeckID != nil {
			wheres = append(wheres, "deck_id=:deck_id")
			arg["deck_id"] = builder.DeckID
		}
		if builder.UserID != nil {
			wheres = append(wheres, "decks.user_id=:user_id")
			arg["user_id"] = builder.UserID
		}
		if builder.Question != nil {
			if builder.Question.Like != nil {
				wheres = append(wheres, "question LIKE :question_like")
				arg["question_like"] = fmt.Sprintf("%%%s%%", *builder.Question.Like)
			}
		}
	}

	if len(wheres) > 0 {
		query = query + " WHERE " + strings.Join(wheres, " AND ")
	}

	return query, arg
}
