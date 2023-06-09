package spec

type Spec struct {
	Services []Service
}

type Service struct {
	Name    string
	Methods []Method
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
	Name           string
	RequestMessage Message
	Kind           CommunicationKind
}

type Message struct {
	Name   string
	Fields []Field
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
	Name string
	Type DataType

	Enum          []string
	IsCollection  bool
	CollectionKey *Field
	OneOf         []Field
	Fields        []Field
}
