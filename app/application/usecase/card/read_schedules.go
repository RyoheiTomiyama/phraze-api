package card

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/domain"
)

func (u *usecase) ReadSchedules(ctx context.Context, cardIDs []int64) (map[int64]*domain.CardSchedule, error) {
	return nil, nil
}
