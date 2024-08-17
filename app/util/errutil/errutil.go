package errutil

import (
	"errors"
	"fmt"
	"runtime"
	"strings"

	"github.com/samber/lo"
)

// サーバエラーが発生した場合にクライアントに返すエラーメッセージ
const InternalErrorMessage = "予期せぬエラーが発生しました"

const MaxStackDepth = 16

type IError interface {
	Error() string
	Message() string
	IsClient() bool
	Code() int
	StackTrace() []uintptr
}

func New(code ErrorCode, format string, args ...interface{}) IError {
	e := fmt.Errorf(format, args...)

	return &customError{
		original: e,
		message:  e.Error(),
		code:     code,
		stack:    caller(0),
	}
}

func Wrap(err error, msg ...string) *customError {
	var ce *customError
	if errors.As(err, &ce) {
		return &customError{
			original: ce,
			message:  lo.Ternary(len(msg) > 0, strings.Join(msg, " "), ce.message),
			code:     ce.code,
			stack:    ce.stack,
		}
	}

	return &customError{
		original: err,
		message:  strings.Join(msg, " "),
		code:     CodeInternalError,
		stack:    caller(0),
	}
}

func As(e error, target interface{}) bool {
	return errors.As(e, target)
}

// stacktraceを取り出す、この関数の位置から積まれてしまうので、skipでうまいこと調整する
func caller(skip int) []uintptr {
	stack := make([]uintptr, MaxStackDepth)
	length := runtime.Callers(3+skip, stack)
	return stack[:length]
}
