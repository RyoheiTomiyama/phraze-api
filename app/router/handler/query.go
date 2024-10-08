package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/errcode"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/RyoheiTomiyama/phraze-api/router/graph/generated"
	"github.com/RyoheiTomiyama/phraze-api/router/graph/resolver"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/RyoheiTomiyama/phraze-api/util/logger"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func PostQuery(r *resolver.Resolver, d *generated.DirectiveRoot) http.HandlerFunc {
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: r, Directives: *d}))

	srv.SetErrorPresenter(func(ctx context.Context, err error) *gqlerror.Error {
		log := logger.FromCtx(ctx)

		defaultError := graphql.DefaultErrorPresenter(ctx, err)

		// NOTE: graphql.DefaultErrorPresenter は引数の err のポインタを返すので、defaultErrPtr のフィールドを変更すると err が書き換わってしまう
		//       元のエラーを変更したくないので値をコピーしてデフォルトのエラーのコードを加工したエラーレスポンスを生成する
		internalErr := *defaultError
		internalErr.Extensions = map[string]interface{}{
			"code":          errutil.CodeInternalError,
			"clientMessage": errutil.InternalErrorMessage,
		}

		var gqlErr *gqlerror.Error
		if ok := errors.As(err, &gqlErr); !ok {
			// err の実体が gqlerror.Error でない場合は FW がバグっている可能性
			log.ErrorWithNotify(ctx, err)

			return &internalErr
		}

		orgErr := gqlErr.Unwrap()
		if orgErr == nil {
			// gqlerror.Error から error を取り出せない場合は FW 内部で生成されたエラー
			codeStr, ok := gqlErr.Extensions["code"].(string)
			if ok {
				// クエリのパースエラー・フィールドのバリデーションエラーはそのまま返す
				if codeStr == errcode.ValidationFailed || codeStr == errcode.ParseFailed {
					log.Error(ctx, err)

					return gqlErr
				}
			} else {
				// parse失敗した場合はバグってる可能性
				log.ErrorWithNotify(ctx, err)
			}

			return &internalErr
		}

		var customErr errutil.IError
		if !errutil.As(orgErr, &customErr) {
			// カスタムエラーにし忘れがある？
			log.ErrorWithNotify(ctx, orgErr, "memo", "カスタムエラーでないものが返されています")

			return &internalErr
		}

		if customErr.IsClient() {
			// クライアントエラーは通知不要
			log.Error(ctx, customErr)

			return &gqlerror.Error{
				Message: customErr.Error(),
				Extensions: map[string]interface{}{
					"code":          customErr.Code(),
					"clientMessage": customErr.Message(),
				},
			}
		}

		log.ErrorWithNotify(ctx, customErr)

		return &internalErr
	})

	return srv.ServeHTTP
}
