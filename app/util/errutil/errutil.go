package errutil

import (
	"errors"
	"fmt"
	"runtime"
	"strings"

	"github.com/samber/lo"
	"golang.org/x/xerrors"
)

// サーバエラーが発生した場合にクライアントに返すエラーメッセージ
const InternalErrorMessage = "予期せぬエラーが発生しました"

const MaxStackDepth = 16

type IError interface {
	Error() string
	Message() string
	// Format(s fmt.State, r rune)
	// FormatError(p xerrors.Printer) error
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
		frame:    xerrors.Caller(1),
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
			frame:    xerrors.Caller(1),
		}
	}

	return &customError{
		original: err,
		message:  strings.Join(msg, " "),
		code:     CodeInternalError,
		stack:    caller(-1),
		frame:    xerrors.Caller(1),
	}
}

func As(e error, target interface{}) bool {
	return errors.As(e, target)
}

func ErrorWithStackTrace(err error) string {
	var ce *customError
	if errors.As(err, &ce) {
		frames := extractFrames(ce.stack)
		traceString := ""
		for _, f := range frames {
			file := f.File
			if strings.HasPrefix(file, "/go/src/app/") {
				relativeIndex := strings.Index(file, "/go/src/app/")
				file = file[relativeIndex+8:]
			}
			function := f.Function
			if i := strings.LastIndex(function, "/"); i > 0 {
				function = function[i+1:]
			}
			traceString += fmt.Sprintf("- %s\n  %s:%d\n", function, file, f.Line)
		}
		return traceString
	}

	return fmt.Sprintf("%+v", err)
}

func caller(skip int) []uintptr {
	stack := make([]uintptr, MaxStackDepth)
	length := runtime.Callers(3+skip, stack)
	return stack[:length]
}

func extractFrames(pcs []uintptr) []runtime.Frame {
	var frames = make([]runtime.Frame, 0, len(pcs))
	callersFrames := runtime.CallersFrames(pcs)

	for {
		callerFrame, more := callersFrames.Next()
		frames = append(frames, callerFrame)

		if !more {
			break
		}
	}

	return frames
}
