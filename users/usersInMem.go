package users

import (
	"errors"
	"sync"

	"github.com/PhilAldridge/TODO-GO/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type InMemoryUsers struct {
	Users
	users []models.User
	mutex sync.RWMutex
}

func NewInMemoryUsersStore() *InMemoryUsers {
	u:= InMemoryUsers{}
	return &u
}

func LoadInMemoryUsersStore(users []models.User) *InMemoryUsers {
	u:=InMemoryUsers{users:users}
	return &u
}

func (u *InMemoryUsers) CreateUser(username string, password string) (uuid.UUID,error) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	newUuid := uuid.New()

	for _,v:= range u.users {
		if v.Username == username {
			return uuid.UUID{}, errors.New("username already exists")
		}
	}

	passwordHash,err:= bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	if err!=nil {
		return uuid.UUID{}, errors.New("unable to secure password")
	}

	u.users = append(u.users, models.User{
		Id: newUuid,
		Username: username,
		PasswordHash: passwordHash,
	})
	return newUuid,nil
}

func (u *InMemoryUsers) Login(username string, password string) (uuid.UUID,error) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	for _,v:= range u.users {
		if v.Username == username {
			err:= bcrypt.CompareHashAndPassword(v.PasswordHash, []byte(password))
			if err != nil {
				return uuid.UUID{}, errors.New("username or password does not match")
			}
			return v.Id,nil
		}
	}

	return uuid.UUID{}, errors.New("username or password does not match")
}
