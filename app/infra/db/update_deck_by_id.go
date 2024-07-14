package db

import (
	"context"
	"time"

	"github.com/RyoheiTomiyama/phraze-api/domain"
)

func (c *client) UpdateCardByID(ctx context.Context, id int64, input *domain.UpdateCardInput) (*domain.Card, error) {
	return &domain.Card{
		CreateAt:  time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
