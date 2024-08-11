package directive

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/RyoheiTomiyama/phraze-api/router/graph/generated"
	"github.com/RyoheiTomiyama/phraze-api/router/graph/model"
	"github.com/RyoheiTomiyama/phraze-api/util/auth"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
)

func New() generated.DirectiveRoot {
	return generated.DirectiveRoot{
		HasRole: HasRole,
	}
}

func HasRole(ctx context.Context, obj interface{}, next graphql.Resolver, role model.Role) (res interface{}, err error) {
	user := auth.FromCtx(ctx)

	//nolint:gocritic
	switch role {
	case model.RoleUser:
		if user == nil {
			return nil, errutil.New(errutil.CodeForbidden, "権限がありません")
		}
	}

	return next(ctx)
}
