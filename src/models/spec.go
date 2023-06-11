package models

type Spec struct {
	Services []Service `json:"services"`
}

type Service struct {
	Name    string   `json:"name"`
	Methods []Method `json:"methods"`
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
	RequestMessage Message           `json:"requestMessage"`
	Kind           CommunicationKind `json:"kind"`
}

type Message struct {
	Name   string  `json:"name"`
	Fields []Field `json:"fields"`
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
	Type         DataType `json:"type"`
	DefaultValue string   `json:"defaultValue"`

	Enum          []string `json:"enum"`
	IsCollection  bool     `json:"isCollection"`
	CollectionKey *Field   `json:"collectionKey"`
	OneOf         []Field  `json:"oneOf"`
	Fields        []Field  `json:"fields"`
}
