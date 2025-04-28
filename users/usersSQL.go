package users

import (
	"database/sql"
	"errors"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type SQLUsers struct {
	Users
	mutex sync.Mutex
	db    *sql.DB
}

func NewSQLUsersStore(db *sql.DB) *SQLUsers {
	u := SQLUsers{db: db}
	return &u
}

func (u *SQLUsers) CreateUser(username string, password string) (uuid.UUID, error) {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	newUUID := uuid.New()
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.UUID{}, errors.New("unable to secure password")
	}

	sqlStatement := `INSERT INTO users (id, username, passwordHash) VALUES ($1,$2,$3)`

	_, err = u.db.Exec(sqlStatement, newUUID, username, passwordHash)
	if err != nil {
		return uuid.UUID{}, err
	}
	return newUUID, nil
}

func (u *SQLUsers) Login(username string, password string) (uuid.UUID, error) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	sqlStatement := `SELECT id, passwordHash FROM users WHERE username = $1`

	res := u.db.QueryRow(sqlStatement, username)
	var userid uuid.UUID
	var passwordHash string
	err := res.Scan(&userid, &passwordHash)
	if err != nil {
		return uuid.UUID{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return uuid.UUID{}, errors.New("username or password does not match")
	}
	return userid, nil
}
