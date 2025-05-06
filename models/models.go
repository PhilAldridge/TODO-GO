package models

import (
	"time"

	"github.com/google/uuid"
)

type Todo struct {
	Id		  uuid.UUID	`json:"Id"`
	Label     string	`json:"Label"`
	AuthorUsername    string	`json:"Author"`
	Completed bool		`json:"Completed"`
	Deadline  time.Time	`json:"Deadline"`
}

type User struct {
	Id				uuid.UUID	`json:"Id"`
	Username		string		`json:"Username"`
	PasswordHash	[]byte		`json:"PasswordHash"`
}

type ContextKey string