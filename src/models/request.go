package models

type Request struct {
	Addr            string `json:"addr"`
	ProtoPath       string `json:"protoPath"`
	FullServiceName string `json:"fullServiceName"`
	MethodName      string `json:"methodName"`
	Script          string `json:"script"`
}

type Response struct {
	Body     string `json:"body"`
	MetaData string `json:"metaData"`
	Logs     string `json:"logs"`
}
