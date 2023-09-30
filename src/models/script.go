package models

import (
	"time"
)

type File struct {
	Name      string    `json:"name"`
	Content   string    `json:"content"`
	Headers   string    `json:"headers"`
	Id 		 string 	  	`json:"id"`
	CreatedAt time.Time `json:"createdAt"  ts_type:"Date" ts_transform:"new Date(__VALUE__)"`
}