package deck

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
)

func (u *usecase) ReadDeckInfo(ctx context.Context, deckIDs []int64) ([]*domain.DeckInfo, error) {
	dmap, err := u.dbClient.GetDeckInfosByDeckID(ctx, deckIDs)
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	// input deckIDsの順番でoutputを作成する
	s := make([]*domain.DeckInfo, len(deckIDs))
	for i, id := range deckIDs {
		s[i] = dmap[id]
	}

	return s, nil
}
