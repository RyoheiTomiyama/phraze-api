package db

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/infra/db/model"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
)

func (c *client) GetPermissionsByUserID(ctx context.Context, userID string) ([]*domain.Permission, error) {
	e := c.execerFrom(ctx)

	query := `
		SELECT 
			permissions.*
		FROM users_roles
			JOIN roles ON roles.id = users_roles.role_id
			JOIN roles_permissions ON roles_permissions.role_id = roles.id
			JOIN permissions ON permissions.id = roles_permissions.permission_id
		WHERE users_roles.user_id = :user_id
		GROUP BY permissions.id
	`
	arg := map[string]interface{}{
		"user_id": userID,
	}

	query, args, err := e.BindNamed(query, arg)
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	var result []*model.Permission
	if err = sqlx.SelectContext(ctx, e, &result, query, args...); err != nil {
		return nil, errutil.Wrap(err)
	}

	permissions := lo.Map(result, func(item *model.Permission, _ int) *domain.Permission {
		return item.ToDomain()
	})

	return permissions, nil
}
