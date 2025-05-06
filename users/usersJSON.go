package users

import (
	"sync"

	"github.com/PhilAldridge/TODO-GO/lib"
	"github.com/google/uuid"
)

type JSONUsers struct {
	Users
	mutex sync.RWMutex
}

func NewJSONUsersStore() *JSONUsers {
	u:= JSONUsers{}
	return &u
}


func (u *JSONUsers) CreateUser(username string, password string) (uuid.UUID,error) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	users:= lib.ReadJsonUsers()
	userStore:=LoadInMemoryUsersStore(users)
	
	uuid,err:= userStore.CreateUser(username, password)
	if err== nil {
		lib.WriteUserStore(userStore.users)
	}

	return uuid,err
}

func (u *JSONUsers) Login(username string, password string) (uuid.UUID,error) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	users:= lib.ReadJsonUsers()
	userStore:=LoadInMemoryUsersStore(users)
	return userStore.Login(username,password)
}

