package auth

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/infra/firebase/auth"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
)

type IAuthUsecase interface {
	ParseToken(ctx context.Context, idToken string) (*domain.User, error)
}

type usecase struct {
	authClient auth.IClient
}

func New(authClient auth.IClient) IAuthUsecase {
	return &usecase{authClient}
}

func (u *usecase) ParseToken(ctx context.Context, idToken string) (*domain.User, error) {
	user, err := u.authClient.Verify(ctx, idToken)
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	return user, nil
}
