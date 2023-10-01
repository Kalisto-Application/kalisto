package models

import (
	"time"
)

type File struct {
	Name      string    `json:"name"`
	Content   string    `json:"content"`
	Headers   string    `json:"headers"`
	Id 		 string 	  	`json:"id"`
}