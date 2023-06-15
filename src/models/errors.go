package models

import "errors"

var (
	ErrProto2NotSupported = errors.New("proto2 is not supported")
	ErrNoProtoFilesFound  = errors.New("no proto files found")
	ErrWorkspaceNotFound  = errors.New("workspace not found")
)
