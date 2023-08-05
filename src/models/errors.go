package models

import "errors"

func NewErrorFormatter(errHandler func(err error)) func(err error) any {
	return func(err error) any {
		var syntaxErr SyntaxError
		if errors.As(err, &syntaxErr) {
			return ApiError{Code: "SYNTAX_ERROR", Value: syntaxErr.Error()}
		}

		errHandler(err)
		return err
	}
}

var (
	ErrProto2NotSupported = errors.New("proto2 is not supported")
	ErrNoProtoFilesFound  = errors.New("no proto files found")
	ErrWorkspaceNotFound  = errors.New("workspace not found")
)

type ApiError struct {
	Code  string
	Value string
}

type SyntaxError string

func (s SyntaxError) Error() string {
	return string(s)
}
