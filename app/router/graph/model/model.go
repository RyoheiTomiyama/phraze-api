package model

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/RyoheiTomiyama/phraze-api/util/logger"
)

// domain→modelに変換する
// dには,*domain.Deckなど、targetには、*model.Deckなどを指定してください
func FromDomain(ctx context.Context, d any, target interface{}) error {
	log := logger.FromCtx(ctx)

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
		t.UserID = v.UserID
		t.Name = v.Name
		t.CreatedAt = v.CreateAt
		t.UpdatedAt = v.UpdatedAt

		return nil
	case *domain.Card:
		if v == nil {
			return nil
		}

		t, ok := target.(*Card)
		if !ok {
			err := errutil.New(errutil.CodeInternalError, "targetとdomainの型が違います")
			log.Error(err, "domain", d, "target", target)

			return err
		}

		t.ID = v.ID
		t.DeckID = v.DeckID
		t.Question = v.Question
		t.Answer = v.Answer
		t.CreatedAt = v.CreateAt
		t.UpdatedAt = v.UpdatedAt

		return nil
	default:
		err := errutil.New(errutil.CodeInternalError, "domain→model変換に失敗しました。")
		log.Error(err, "domain", d)

		return err
	}
}
