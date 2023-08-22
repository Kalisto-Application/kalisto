package models

type ProtoDir struct {
	Folder string   `json:"folder"`
	Files  []string `json:"files"`
}
