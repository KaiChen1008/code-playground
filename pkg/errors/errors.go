package errors

import (
	stderrors "errors"
	"fmt"
)

type Error struct {
	Message string `json:"message"`
	Err     error  `json:"inner_error,omitempty"`
}

func New(args ...any) error {
	// Keep the behavior of standard errors package
	if len(args) == 1 {
		if text, ok := args[0].(string); ok {
			return stderrors.New(text)
		}
	}

	ret := Error{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case *Error:
			ret.setErr(arg.Err)
			ret.appendMsg(arg.Message)
		case error:
			ret.setErr(arg)
			ret.appendMsg(arg.Error())
		case string:
			ret.appendMsg(arg)
		default:
			ret.appendMsg(fmt.Sprintf("%v", arg))
		}
	}
	return &ret
}

func (e *Error) setErr(err error) {
	if e.Err == nil {
		e.Err = err
	}
}
func (e *Error) appendMsg(msg string) {
	if e.Message == "" {
		e.Message = msg
		return
	}
	e.Message = fmt.Sprintf("%s: %s", e.Message, msg)
}

// Error returns an human-readable error message.
func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

func Is(err, target error) bool {
	return stderrors.Is(err, target)
}

func As(err error, target any) bool {
	return stderrors.As(err, target)
}

func Unwrap(err error) error {
	return stderrors.Unwrap(err)
}
