package auth

import (
	"context"
	"encoding/json"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
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

type customClaims struct {
	Name    string `json:"name"`
	Picture string `json:"picture"`
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
}

func (c *client) Verify(ctx context.Context, idToken string) (*domain.User, error) {
	// VerifyIDTokenAndCheckRevoked は遅いので使わない
	token, err := c.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	jsonstr, err := json.Marshal(token.Claims)
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	var claims customClaims
	err = json.Unmarshal(jsonstr, &claims)
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	return &domain.User{
		ID:     claims.UserID,
		Name:   &claims.Name,
		Avatar: &claims.Picture,
	}, nil
}
