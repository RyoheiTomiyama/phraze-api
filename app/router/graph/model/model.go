package model

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/RyoheiTomiyama/phraze-api/util/logger"
)

func FromDomain(ctx context.Context, d any, target interface{}) error {
	log := logger.FromCtx(ctx)

	log.Debug("Fromdomai", "d", d, "target", target)

	switch v := d.(type) {

	case *domain.Deck:
		if v == nil {
			return nil
		}

		t, ok := target.(*Deck)
		if !ok {
			err := errutil.New(errutil.CodeInternalError, "targetとdomainの型が違います")
			log.Error(err, "domain", d, "target", target)

			return err
		}

		t.ID = v.ID
		t.Name = v.Name
		t.CreatedAt = v.CreateAt
		t.UpdatedAt = v.UpdatedAt

		return nil
	default:
		err := errutil.New(errutil.CodeInternalError, "domain→model変換に失敗しました。")
		log.Error(err, "domain", d)

		return err
	}
}
