package errutil

type customError struct {
	code     ErrorCode
	message  string
	original error
	stack    []uintptr
}

// debugエラーメッセージ
func (c *customError) Error() string {
	return c.original.Error()
}

// クライアントに返すエラーメッセージ
func (c *customError) Message() string {
	if c.code.IsClientError() {
		return c.message
	}

	return InternalErrorMessage
}

func (c *customError) IsClient() bool {
	return c.code.IsClientError()
}

func (c *customError) Code() int {
	return int(c.code)
}

// SentryにStackTraceを表示するために、このメソッドが必要
func (c *customError) StackTrace() []uintptr {
	return c.stack
}
