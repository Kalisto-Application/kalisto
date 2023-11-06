package models

import (
	"context"
	"errors"
	"fmt"
	"kalisto/src/pkg/runtime"
	"log"

	rpkg "github.com/wailsapp/wails/v2/pkg/runtime"
)

func NewErrorFormatter(ctxGetter func() context.Context, errHandler func(err error), runtime runtime.Runtime) func(err error) any {
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
		/////
		var errFileMustBeAbsolute *ErrorFileMustBeAbsolute
		if errors.As(err, &errFileMustBeAbsolute) {
			runtime.MessageDialog(ctxGetter(), rpkg.MessageDialogOptions{
				Type:    "error",
				Title:   "File must be absolute",
				Message: errFileMustBeAbsolute.File,
			})
			return ApiError{Code: "FILE_MUST_BE_ABSOLUTE", Value: errFileMustBeAbsolute.File}
		}
		/////
		var errOpenapiFileCantBeDir *ErrorOpenapiFileCantBeDir
		if errors.As(err, &errOpenapiFileCantBeDir) {
			runtime.MessageDialog(ctxGetter(), rpkg.MessageDialogOptions{
				Type:    "error",
				Title:   "File can't be a directory",
				Message: errOpenapiFileCantBeDir.File,
			})
			return ApiError{Code: "OPENAPI_FILE_CANT_BE_DIR", Value: errFileMustBeAbsolute.File}
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
	ErrMethodNotFound     = errors.New("method not found")
)

type ApiError struct {
	Code  string
	Value string
}

type ErrorSyntax string

func (s ErrorSyntax) Error() string {
	return string(s)
}

type ErrorFileMustBeAbsolute struct {
	File string
}

func (e *ErrorFileMustBeAbsolute) Error() string {
	return "filename must be absolute"
}

type ErrorOpenapiFileCantBeDir struct {
	File string
}

func (e *ErrorOpenapiFileCantBeDir) Error() string {
	return "openapi file can't be a directory"
}

var JsTypeError = fmt.Errorf("TypeError: expected an object")

var ErrorServerUnavailable = errors.New("server unavailable")
