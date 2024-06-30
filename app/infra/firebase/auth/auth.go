package auth

import (
	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/RyoheiTomiyama/phraze-api/util/logger"
)

type client struct {
	client *auth.Client
}
type IClient interface {
	Verify(ctx context.Context, idToken string) (*domain.User, error)
}

func New() (IClient, error) {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	// Access auth service from the default app
	c, err := app.Auth(context.Background())
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	return &client{
		client: c,
	}, nil
}

func (c *client) Verify(ctx context.Context, idToken string) (*domain.User, error) {
	l := logger.FromCtx(ctx)

	token, err := c.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	l.Debug("firebase auth verify token", "token", token)

	return &domain.User{
		ID: token.UID,
	}, nil
}
