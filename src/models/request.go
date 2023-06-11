package models

type Request struct {
	Body     string `json:"body"`
	MetaData string `json:"metaData"`
}

type Response struct {
	Body     string `json:"body"`
	MetaData string `json:"metaData"`
	Logs     string `json:"logs"`
}
