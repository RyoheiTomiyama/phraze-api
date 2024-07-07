package handler

import (
	"context"
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

func PostQuery(r *resolver.Resolver) http.HandlerFunc {
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: r}))

	srv.SetErrorPresenter(func(ctx context.Context, err error) *gqlerror.Error {
		log := logger.FromCtx(ctx)

		defaultError := graphql.DefaultErrorPresenter(ctx, err)

		// NOTE: graphql.DefaultErrorPresenter は引数の err のポインタを返すので、defaultErrPtr のフィールドを変更すると err が書き換わってしまう
		//       元のエラーを変更したくないので値をコピーしてデフォルトのエラーのコードを加工したエラーレスポンスを生成する
		internalErr := *defaultError
		internalErr.Extensions = map[string]interface{}{
			"code": errutil.CodeInternalError,
		}

		gqlErr, ok := err.(*gqlerror.Error)
		if !ok {
			// err の実体が gqlerror.Error でない場合は FW がバグっている可能性
			log.Error(err)
			// TODO Sentry

			return &internalErr
		}

		orgErr := gqlErr.Unwrap()
		if orgErr == nil {
			// gqlerror.Error から error を取り出せない場合は FW 内部で生成されたエラー
			codeStr, ok := gqlErr.Extensions["code"].(string)
			if ok {
				// クエリのパースエラー・フィールドのバリデーションエラーはそのまま返す
				if codeStr == errcode.ValidationFailed || codeStr == errcode.ParseFailed {
					log.Error(err)

					return gqlErr
				}
			} else {
				// parse失敗した場合はバグってる可能性
				log.Error(err)
			}

			return &internalErr
		}

		var customErr errutil.IError
		if !errutil.As(orgErr, &customErr) {
			// カスタムエラーにし忘れがある？
			log.Error(orgErr, "memo", "カスタムエラーでないものが返されています")
			// TODO Sentry

			return &internalErr
		}

		log.Error(customErr)
		if customErr.IsClient() {
			return &gqlerror.Error{
				Message: customErr.Error(),
				Extensions: map[string]interface{}{
					"code":          customErr.Code(),
					"clientMessage": customErr.Message(),
				},
			}
		}

		// TODO Sentry
		return &internalErr
	})

	return srv.ServeHTTP
}
