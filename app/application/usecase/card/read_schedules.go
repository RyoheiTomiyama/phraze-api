package card

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
)

func (u *usecase) ReadSchedules(ctx context.Context, cardIDs []int64) (map[int64]*domain.CardSchedule, error) {
	smap, err := u.dbClient.GetCardSchedulesByCardID(ctx, cardIDs)
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	return smap, nil
}
