package models

import (
	"time"

	"github.com/google/uuid"
)

type Todo struct {
	Id		  uuid.UUID	`json:"Id"`
	Label     string	`json:"Label"`
	Author    string	`json:"Author"`
	Completed bool		`json:"Completed"`
	Deadline  time.Time	`json:"Deadline"`
}

