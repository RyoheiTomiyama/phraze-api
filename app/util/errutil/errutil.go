package errutil

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/xerrors"
)

// サーバエラーが発生した場合にクライアントに返すエラーメッセージ
const InternalErrorMessage = "予期せぬエラーが発生しました"

func New(code errorCode, format string, args ...interface{}) *customError {
	e := fmt.Errorf(format, args...)

	return &customError{original: e, message: e.Error(), code: code, frame: xerrors.Caller(1)}
}

func Wrap(err error, msg ...string) *customError {
	var ce *customError
	if errors.As(err, &ce) {
		return &customError{
			original: ce,
			message:  strings.Join(msg, " "),
			code:     ce.code,
			frame:    xerrors.Caller(1),
		}
	}

	return &customError{
		original: err,
		message:  strings.Join(msg, " "),
		code:     CodeInternalError,
		frame:    xerrors.Caller(1),
	}
}

func ErrorWithStackTrace(err error) string {
	return fmt.Sprintf("%+v", err)
}
