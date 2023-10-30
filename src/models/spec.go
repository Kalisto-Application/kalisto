package models

import (
	"fmt"
	"strings"
	"time"
)

type WorkspaceList struct {
	Main Workspace        `json:"main"`
	List []WorkspaceShort `json:"list"`
}

type WorkspaceShort struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Workspace struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	TargetUrl    string    `json:"targetUrl"`
	Spec         Spec      `json:"spec"`
	LastUsage    time.Time `json:"lastUsage" ts_type:"Date" ts_transform:"new Date(__VALUE__)"`
	BasePath     []string  `json:"basePath"`
	RequestFiles []File    `json:"requestFiles"`
	ScriptFiles  []File    `json:"scriptFiles"`
}

type WorkspaceKind string

const (
	WorkspaceKindProto   WorkspaceKind = "proto"
	WorkspaceKindOpenapi WorkspaceKind = "openapi"
)

type Spec struct {
	Services []Service          `json:"services"`
	Links    map[string]Message `json:"-"`
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

func (s *Spec) FindOutputMessage(serviceName, methodName string) (Message, error) {
	for _, service := range s.Services {
		if service.FullName == serviceName {
			for _, method := range service.Methods {
				if method.Name == methodName {
					return method.ResponseMessage, nil
				}
			}
		}
	}

	return Message{}, fmt.Errorf("output message not found, method=%s", methodName)
}

type Service struct {
	Name        string   `json:"name"`
	DisplayName string   `json:"displayName"`
	Package     string   `json:"package"`
	FullName    string   `json:"fullName"`
	Methods     []Method `json:"methods"`
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
	Name             string            `json:"name"`
	FullName         string            `json:"fullName"`
	RequestMessage   Message           `json:"requestMessage"`
	ResponseMessage  Message           `json:"responseMessage"`
	Kind             CommunicationKind `json:"kind"`
	RequestExample   string            `json:"requestExample"`
	RequestInstances []File            `json:"requestInstances"`
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

	DataTypeDate     DataType = "DataTypeDate"
	DataTypeDuration DataType = "DataTypeDuration"

	DataTypeOneOf DataType = "DataTypeOneOf"
)

type Field struct {
	Name     string   `json:"name"`
	FullName string   `json:"fullName"`
	Type     DataType `json:"type"`

	Enum     []int32 `json:"enum"`
	Repeated bool    `json:"repeated"`
	MapKey   *Field  `json:"mapKey"`
	MapValue *Field  `json:"mapValue"`
	OneOf    []Field `json:"oneOf"`
	Message  string  `json:"message"`
}

func (f Field) FindOneofByName(name string) (Field, error) {
	for i := range f.OneOf {
		if f.OneOf[i].Name == name {
			return f.OneOf[i], nil
		}
	}

	return Field{}, fmt.Errorf("can't find a field<%s>", name)
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
