package models

import "time"

type File struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Content   string    `json:"content"`
	Headers   string    `json:"headers"`
	CreatedAt time.Time `json:"createdAt" ts_type:"Date" ts_transform:"new Date(__VALUE__)"`
}
