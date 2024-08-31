package db

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/domain"
)

func (c *client) HasPermissionByUserID(ctx context.Context, userID string, key domain.PermissionKey) (bool, error) {
	return false, nil
}
