package errutil

type errorCode int

const (
	// クライアントエラー系
	CodeBadRequest   errorCode = 400
	CodeUnauthorized errorCode = 401
	CodeForbidden    errorCode = 403
	CodeNotFound     errorCode = 404
	CodeConflict     errorCode = 407

	// サーバエラー系
	CodeInternalError errorCode = 500
)

func (c errorCode) IsClientError() bool {
	return c < 500
}
