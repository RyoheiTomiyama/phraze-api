package errutil

type ErrorCode int

const (
	// クライアントエラー系
	CodeBadRequest   ErrorCode = 400
	CodeUnauthorized ErrorCode = 401
	CodeForbidden    ErrorCode = 403
	CodeNotFound     ErrorCode = 404
	CodeConflict     ErrorCode = 407

	// サーバエラー系
	CodeInternalError ErrorCode = 500
)

func (c ErrorCode) IsClientError() bool {
	return c < 500
}
