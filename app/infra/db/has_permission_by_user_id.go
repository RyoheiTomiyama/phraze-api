package db

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/jmoiron/sqlx"
)

func (c *client) HasPermissionByUserID(ctx context.Context, userID string, key domain.PermissionKey) (bool, error) {
	e := c.execerFrom(ctx)

	query := `
		SELECT 
			COUNT(users_roles.user_id) > 0
		FROM users_roles
			JOIN roles ON roles.id = users_roles.role_id
			JOIN roles_permissions ON roles_permissions.role_id = roles.id
			JOIN permissions ON permissions.id = roles_permissions.permission_id
		WHERE users_roles.user_id = :user_id AND permissions.key = :key
	`
	arg := map[string]interface{}{
		"user_id": userID,
		"key":     key,
	}

	query, args, err := e.BindNamed(query, arg)
	if err != nil {
		return false, errutil.Wrap(err)
	}

	var result bool
	if err = sqlx.GetContext(ctx, e, &result, query, args...); err != nil {
		return false, errutil.Wrap(err)
	}

	return result, nil
}
