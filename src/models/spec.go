package models

import (
	"time"

	"github.com/google/uuid"
)

type Workspace struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"workspace"`
	Spec      Spec      `json:"spec"`
	LastUsage time.Time `json:"lastUsage"`
	BasePath  string    `json:"basePath"`
}

type Spec struct {
	Services []Service `json:"services"`
}
type Service struct {
	Name     string   `json:"name"`
	Package  string   `json:"package"`
	FullName string   `json:"fullName"`
	Methods  []Method `json:"methods"`
}

type CommunicationKind string

const (
	CommunicationKindSimple       CommunicationKind = "Simple"
	CommunicationKindClientStream CommunicationKind = "ClientStream"
	CommunicationKindServerStream CommunicationKind = "ServerStream"
	CommunicationKindBidirection  CommunicationKind = "Bidirection"
)

func NewCommunicationKind(isStreamClient, isStreamServer bool) CommunicationKind {
	if isStreamClient && isStreamServer {
		return CommunicationKindBidirection
	}
	if isStreamClient {
		return CommunicationKindClientStream
	}
	if isStreamServer {
		return CommunicationKindServerStream
	}
	return CommunicationKindSimple
}

type Method struct {
	Name           string            `json:"name"`
	FullName       string            `json:"fullName"`
	RequestMessage Message           `json:"requestMessage"`
	Kind           CommunicationKind `json:"kind"`
}

type Message struct {
	Name     string  `json:"name"`
	FullName string  `json:"fullName"`
	Fields   []Field `json:"fields"`
}

type DataType string

const (
	DataTypeBool   DataType = "DataTypeBool"
	DataTypeInt    DataType = "DataTypeInt"
	DataTypeFloat  DataType = "DataTypeFloat"
	DataTypeString DataType = "DataTypeString"
	DataTypeEnum   DataType = "DataTypeEnum"

	DataTypeStruct DataType = "DataTypeStruct"
	DataTypeMap    DataType = "DataTypeMap"
	DataTypeArray  DataType = "DataTypeArray"

	DataTypeDate DataType = "DataTypeDate"
)

type Field struct {
	Name         string   `json:"name"`
	FullName     string   `json:"fullName"`
	Type         DataType `json:"type"`
	DefaultValue string   `json:"defaultValue"`

	Enum          []string `json:"enum"`
	IsCollection  bool     `json:"isCollection"`
	CollectionKey *Field   `json:"collectionKey"`
	OneOf         []Field  `json:"oneOf"`
	Fields        []Field  `json:"fields"`
}
