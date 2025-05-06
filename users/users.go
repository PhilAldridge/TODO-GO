package users

import (
	"github.com/google/uuid"
)

type Users interface {
	//GetUser(id uuid.UUID) (models.User,error)
	CreateUser(username string, password string) (uuid.UUID,error)
	//DeleteUser(id uuid.UUID) error
	Login(username string, password string) (uuid.UUID,error)
}