package errutil

import (
	"errors"
	"fmt"
	"strings"

	"github.com/samber/lo"
	"golang.org/x/xerrors"
)

// サーバエラーが発生した場合にクライアントに返すエラーメッセージ
const InternalErrorMessage = "予期せぬエラーが発生しました"

type IError interface {
	Error() string
	Message() string
	Format(s fmt.State, r rune)
	FormatError(p xerrors.Printer) error
	IsClient() bool
	Code() int
}

func New(code ErrorCode, format string, args ...interface{}) IError {
	e := fmt.Errorf(format, args...)

	return &customError{original: e, message: e.Error(), code: code, frame: xerrors.Caller(1)}
}

func Wrap(err error, msg ...string) *customError {
	var ce *customError
	if errors.As(err, &ce) {
		return &customError{
			original: ce,
			message:  lo.Ternary(len(msg) > 0, strings.Join(msg, " "), ce.message),
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

func As(e error, target interface{}) bool {
	return errors.As(e, target)
}

func ErrorWithStackTrace(err error) string {
	return fmt.Sprintf("%+v", err)
}
