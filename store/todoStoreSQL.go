package store

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/PhilAldridge/TODO-GO/lib"
	"github.com/PhilAldridge/TODO-GO/models"
	"github.com/google/uuid"
)

type SQLTodoStore struct {
	Store
	mutex sync.Mutex
	db    *sql.DB
}

func NewSQLStore() *SQLTodoStore {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		lib.SqlHost, lib.SqlPortNo, lib.SqlUser, lib.SqlPassword, lib.SqlDbName)

	fmt.Println(psqlInfo)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	newStore := &SQLTodoStore{db: db}
	newStore.setupDB()
	return newStore
}

func (t *SQLTodoStore) Close() {
	t.db.Close()
}

func (t *SQLTodoStore) AddTodo(label string,deadline time.Time,	username string) (uuid.UUID, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	newUuid:= uuid.New()
	sqlStatement:= `INSERT INTO todos (id, label, deadline, completed, authorusername)
	VALUES ($1, $2, $3, false, $4)`
	_,err:= t.db.Exec(sqlStatement, newUuid,label,deadline,username)
	if err!=nil {
		return uuid.UUID{}, err
	}
	return newUuid,nil
}

func (t *SQLTodoStore) GetTodos(username string) []models.Todo {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	sqlStatement:= `SELECT 
		id
		,label
		,authorusername
		,deadline
		,completed
	FROM todos WHERE authorusername = $1`
	res, _:= t.db.Query(sqlStatement, username)
	defer res.Close()

	var todos []models.Todo
	for res.Next() {
		todo:= models.Todo{}
		res.Scan(
			&todo.Id,
			&todo.Label, 
			&todo.AuthorUsername,
			&todo.Deadline,
			&todo.Completed,
		)
		todos = append(todos, todo)
	}
	return todos
}

func (t *SQLTodoStore) GetTodoById(id uuid.UUID, username string) (models.Todo, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	sqlStatement:= `SELECT 
		id
		,label
		,authorusername
		,deadline
		,completed
	FROM todos WHERE authorusername = $1 AND id = $2`
	res:= t.db.QueryRow(sqlStatement, username, id)

	todo:= models.Todo{}
	switch err:= res.Scan(
		&todo.Id,
		&todo.Label, 
		&todo.AuthorUsername,
		&todo.Deadline,
		&todo.Completed,
	); err {
	case sql.ErrNoRows:
		return models.Todo{},errors.New("no matching todo found")
	case nil:
		return todo,nil
	default:
		return models.Todo{},err
	}	
}

func (t *SQLTodoStore) UpdateTodo(id uuid.UUID,
	field string,
	updatedValue string,
	username string,
) (models.Todo, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	sqlStatement:= fmt.Sprintf(`UPDATE todos
	SET %s = $1
	WHERE id = $2 AND authorusername = $3
	RETURNING 
		id
		,label
		,authorusername
		,deadline
		,completed`,field)
	res:= t.db.QueryRow(sqlStatement, updatedValue, id, username)
	todo:= models.Todo{}
	switch err:= res.Scan(
		&todo.Id,
		&todo.Label, 
		&todo.AuthorUsername,
		&todo.Deadline,
		&todo.Completed,
	); err {
	case sql.ErrNoRows:
		return models.Todo{},errors.New("no matching todo found")
	case nil:
		return todo,nil
	default:
		return models.Todo{},err
	}	
}

func (t *SQLTodoStore) DeleteTodo(id uuid.UUID,	username string)  error {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	sqlStatement:= `DELETE FROM todos
	WHERE id = $1 AND authorusername = $2`
	res,err:= t.db.Exec(sqlStatement, id, username)
	if err!=nil {
		return err
	}
	if count,err:= res.RowsAffected(); count==0 || err!=nil {
		return errors.New("todo not found")
	}
	return nil
}


func (t *SQLTodoStore) setupDB() {
	sqlStatement := `CREATE TABLE IF NOT EXISTS users (
		id uuid Primary Key,
		username TEXT
	)`
	_, err := t.db.Exec(sqlStatement)
	if err != nil {
		panic(err)
	}

	sqlStatement = `CREATE TABLE IF NOT EXISTS todos (
		id uuid primary key,
		label     text,
		authorusername    text,
		completed boolean,
		deadline  timestamp
	)`
	_, err = t.db.Exec(sqlStatement)
	if err != nil {
		panic(err)
	}
}
