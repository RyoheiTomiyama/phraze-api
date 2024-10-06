package firebase

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/domain"
)

type IClient interface {
	Verify(ctx context.Context, idToken string) (*domain.User, error)
}
