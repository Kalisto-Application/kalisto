package models

import (
	"errors"
	"fmt"
	"log"
)

func NewErrorFormatter(errHandler func(err error)) func(err error) any {
	return func(err error) any {
		var syntaxErr ErrorSyntax
		if errors.As(err, &syntaxErr) {
			return ApiError{Code: "SYNTAX_ERROR", Value: syntaxErr.Error()}
		}
		if errors.Is(err, ErrorServerUnavailable) {
			return ApiError{Code: "SERVER_UNAVAILABLE", Value: ""}
		}
		if errors.Is(err, JsTypeError) {
			return ApiError{Code: "SYNTAX_ERROR", Value: err.Error()}
		}

		errHandler(err)
		log.Println("api error: ", err.Error())
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

type ErrorSyntax string

func (s ErrorSyntax) Error() string {
	return string(s)
}

var JsTypeError = fmt.Errorf("TypeError: expected an object")

var ErrorServerUnavailable = errors.New("server unavailable")
