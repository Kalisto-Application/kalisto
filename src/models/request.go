package models

type Request struct {
	Addr        string `json:"addr"`
	WorkspaceID string `json:"workspaceId"`
	Method      string `json:"method"`
	Body        string `json:"body"`
	Meta        string `json:"meta"`
}

type Response struct {
	Body     string `json:"body"`
	MetaData string `json:"metaData"`
	Logs     string `json:"logs"`
}

type ScriptCall struct {
	Addr        string `json:"addr"`
	WorkspaceID string `json:"workspaceId"`
	Body        string `json:"body"`
	Meta        string `json:"meta"`
}
