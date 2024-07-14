package db

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/domain"
)

func (c *client) GetCards(ctx context.Context, where *domain.CardsWhere, limit, offset int) ([]*domain.Card, error) {
	return nil, nil
}
