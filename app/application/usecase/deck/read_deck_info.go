package deck

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/domain"
)

func (u *usecase) ReadDeckInfo(ctx context.Context, deckIDs []int64) ([]*domain.DeckInfo, error) {
	return nil, nil
}
