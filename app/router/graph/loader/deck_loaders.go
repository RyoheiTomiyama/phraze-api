package loader

import (
	"context"
	"time"

	"github.com/RyoheiTomiyama/phraze-api/application/usecase/card"
	"github.com/vikstrous/dataloadgen"
)

type deckLoader struct {
	ScheduleAtLoader LoaderInterface[int64, time.Time]
}

type IDeckLoader interface {
}

func NewDeckLoader(cardUsecase card.IUsecase) IDeckLoader {
	reader := &deckReader{cardUsecase}

	return &deckLoader{
		ScheduleAtLoader: NewNoCacheLoader(reader.ReadScheduleAt, dataloadgen.WithWait(time.Millisecond)),
	}
}

type deckReader struct {
	cardUsecase card.IUsecase
}

func (r *deckReader) ReadScheduleAt(ctx context.Context, deckIDs []int64) ([]time.Time, []error) {
	return nil, nil
}
