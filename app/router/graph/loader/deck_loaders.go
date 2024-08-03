package loader

import (
	"context"
	"time"

	"github.com/RyoheiTomiyama/phraze-api/application/usecase/deck"
	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/vikstrous/dataloadgen"
)

type deckLoader struct {
	DeckInfoLoader LoaderInterface[int64, *domain.DeckInfo]
}

type IDeckLoader interface {
	GetScheduleAt(ctx context.Context, deckID int64) (*domain.DeckInfo, error)
}

func NewDeckLoader(deckUsecase deck.IUsecase) IDeckLoader {
	reader := &deckReader{deckUsecase}

	return &deckLoader{
		DeckInfoLoader: NewNoCacheLoader(reader.ReadDeckInfo, dataloadgen.WithWait(time.Millisecond)),
	}
}

func (l *deckLoader) GetScheduleAt(ctx context.Context, deckID int64) (*domain.DeckInfo, error) {
	sa, err := l.DeckInfoLoader.Load(ctx, deckID)
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	return sa, nil
}

type deckReader struct {
	deckUsecase deck.IUsecase
}

func (r *deckReader) ReadDeckInfo(ctx context.Context, deckIDs []int64) ([]*domain.DeckInfo, []error) {
	deckInfoList, err := r.deckUsecase.ReadDeckInfo(ctx, deckIDs)
	if err != nil {
		return nil, []error{errutil.Wrap(err)}
	}

	return deckInfoList, nil
}
