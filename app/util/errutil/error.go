package errutil

import (
	"fmt"

	"golang.org/x/xerrors"
)

type customError struct {
	code     errorCode
	message  string
	original error
	frame    xerrors.Frame
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

func (c *customError) Format(s fmt.State, r rune) { // implements fmt.Formatter
	xerrors.FormatError(c, s, r)
}

func (c *customError) FormatError(p xerrors.Printer) error { // implements xerrors.Formatter
	p.Print(c.message)
	if p.Detail() {
		c.frame.Format(p)
	}

	return c.original
}

func (c *customError) IsClient() bool {
	return c.code.IsClientError()
}

func (c *customError) Code() int {
	return int(c.code)
}
