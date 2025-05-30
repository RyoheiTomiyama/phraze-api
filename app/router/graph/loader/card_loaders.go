package loader

import (
	"context"
	"time"

	"github.com/RyoheiTomiyama/phraze-api/application/usecase/card"
	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/vikstrous/dataloadgen"
)

type cardLoader struct {
	ScheduleLoader LoaderInterface[int64, *domain.CardSchedule]
}

type ICardLoader interface {
	GetSchedule(ctx context.Context, cardID int64) (*domain.CardSchedule, error)
}

func NewCardLoader(cardUsecase card.IUsecase) ICardLoader {
	reader := &cardReader{cardUsecase}

	return &cardLoader{
		ScheduleLoader: NewNoCacheLoader(reader.ReadSchedules, dataloadgen.WithWait(time.Millisecond)),
	}
}

func (l *cardLoader) GetSchedule(ctx context.Context, cardID int64) (*domain.CardSchedule, error) {
	s, err := l.ScheduleLoader.Load(ctx, cardID)
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	return s, nil
}

type cardReader struct {
	cardUsecase card.IUsecase
}

func (r *cardReader) ReadSchedules(ctx context.Context, cardIDs []int64) ([]*domain.CardSchedule, []error) {
	scheduleMap, err := r.cardUsecase.ReadSchedules(ctx, cardIDs)
	if err != nil {
		return nil, []error{errutil.Wrap(err)}
	}

	schedules := make([]*domain.CardSchedule, len(cardIDs))

	for i, id := range cardIDs {
		schedules[i] = scheduleMap[id]
	}

	return schedules, nil
}
