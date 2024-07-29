package loader

import (
	"context"
	"time"

	"github.com/RyoheiTomiyama/phraze-api/application/usecase/card"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/vikstrous/dataloadgen"
)

type deckLoader struct {
	ScheduleAtLoader LoaderInterface[int64, *time.Time]
}

type IDeckLoader interface {
	GetScheduleAt(ctx context.Context, deckID int64) (*time.Time, error)
}

func NewDeckLoader(cardUsecase card.IUsecase) IDeckLoader {
	reader := &deckReader{cardUsecase}

	return &deckLoader{
		ScheduleAtLoader: NewNoCacheLoader(reader.ReadScheduleAt, dataloadgen.WithWait(time.Millisecond)),
	}
}

func (l *deckLoader) GetScheduleAt(ctx context.Context, deckID int64) (*time.Time, error) {
	sa, err := l.ScheduleAtLoader.Load(ctx, deckID)
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	return sa, nil
}

type deckReader struct {
	cardUsecase card.IUsecase
}

func (r *deckReader) ReadScheduleAt(ctx context.Context, deckIDs []int64) ([]*time.Time, []error) {
	schedules, err := r.cardUsecase.ReadScheduleAt(ctx, deckIDs)
	if err != nil {
		return nil, []error{errutil.Wrap(err)}
	}

	return schedules, nil
}
