package models

import (
	"fmt"
	"strings"
	"time"
)

type Workspace struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Spec      Spec      `json:"spec"`
	LastUsage time.Time `json:"lastUsage"`
	BasePath  string    `json:"basePath"`
}

type Spec struct {
	Services []Service          `json:"services"`
	Links    map[string]Message `json:"links"`
}

func (s *Spec) FindInputMessage(serviceName, methodName string) (Message, error) {
	for _, service := range s.Services {
		if service.FullName == serviceName {
			for _, method := range service.Methods {
				if method.Name == methodName {
					return method.RequestMessage, nil
				}
			}
		}
	}

	return Message{}, fmt.Errorf("input message not found, method=%s", methodName)
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
	RequestExample string            `json:"requestExample"`
}

type Message struct {
	Name     string  `json:"name"`
	FullName string  `json:"fullName"`
	Fields   []Field `json:"fields"`
}

func (m Message) FindField(name string) (Field, error) {
	for _, f := range m.Fields {
		if f.Name == name {
			return f, nil
		}
	}

	return Field{}, fmt.Errorf("field=%s not found", name)
}

type DataType string

const (
	DataTypeBool    DataType = "DataTypeBool"
	DataTypeInt32   DataType = "DataTypeInt32"
	DataTypeInt64   DataType = "DataTypeInt64"
	DataTypeUint32  DataType = "DataTypeUint32"
	DataTypeUint64  DataType = "DataTypeUint64"
	DataTypeFloat32 DataType = "DataTypeFloat32"
	DataTypeFloat64 DataType = "DataTypeFloat64"
	DataTypeString  DataType = "DataTypeString"
	DataTypeBytes   DataType = "DataTypeBytes"

	DataTypeEnum DataType = "DataTypeEnum"

	DataTypeStruct DataType = "DataTypeStruct"
	DataTypeMap    DataType = "DataTypeMap"
	DataTypeArray  DataType = "DataTypeArray"

	DataTypeDate DataType = "DataTypeDate"
)

type CastType string

const (
	CastTypeBool    CastType = "CastTypeBool"
	CastTypeInt32   CastType = "CastTypeInt32"
	CastTypeInt64   CastType = "CastTypeInt64"
	CastTypeUint32  CastType = "CastTypeUint32"
	CastTypeUint64  CastType = "CastTypeUint64"
	CastTypeFloat32 CastType = "CastTypeFloat32"
	CastTypeFloat64 CastType = "CastTypeFloat64"
	CastTypeString  CastType = "CastTypeString"
	CastTypeBytes   CastType = "CastTypeBytes"

	CastTypeEnum CastType = "CastTypeEnum"

	CastTypeStruct CastType = "CastTypeStruct"
	CastTypeMap    CastType = "CastTypeMap"
	CastTypeArray  CastType = "CastTypeArray"

	CastTypeDate CastType = "CastTypeDate"
)

var DataToCast = map[DataType]CastType{
	DataTypeBool:    CastTypeBool,
	DataTypeInt32:   CastTypeInt32,
	DataTypeInt64:   CastTypeInt64,
	DataTypeUint32:  CastTypeUint32,
	DataTypeUint64:  CastTypeUint64,
	DataTypeFloat32: CastTypeFloat32,
	DataTypeFloat64: CastTypeFloat64,
	DataTypeString:  CastTypeString,
	DataTypeBytes:   CastTypeBytes,
	DataTypeEnum:    CastTypeEnum,
	DataTypeStruct:  CastTypeStruct,
	DataTypeMap:     CastTypeMap,
	DataTypeArray:   CastTypeArray,
	DataTypeDate:    CastTypeDate,
}

type Field struct {
	Name         string   `json:"name"`
	FullName     string   `json:"fullName"`
	Type         DataType `json:"type"`
	DefaultValue string   `json:"defaultValue"`

	Enum          []int32 `json:"enum"`
	IsCollection  bool    `json:"isCollection"`
	CollectionKey *Field  `json:"collectionKey"`
	OneOf         []Field `json:"oneOf"`
	Message       string  `json:"message"`
}

type MethodName string

func (m MethodName) ServiceAndShort() (string, string) {
	values := strings.Split(string(m), ".")
	if len(values) == 0 {
		return "", ""
	}

	s := values[0 : len(values)-1]
	return strings.Join(s, "."), values[len(values)-1]
}
