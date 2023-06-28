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
				if method.FullName == methodName {
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
	Message       string   `json:"message"`
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
