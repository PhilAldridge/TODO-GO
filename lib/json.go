package lib

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/PhilAldridge/TODO-GO/models"
)

func ReadJsonStore() []models.Todo {
	file, err := os.Open(JsonStoreFile)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var todos []models.Todo
	byteValue, err := io.ReadAll(file)

	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(byteValue, &todos)

	return todos
}

func WriteJsonStore(data []models.Todo) {
	jsonString, _ := json.Marshal(data)
	os.WriteFile(JsonStoreFile, jsonString, os.ModePerm)
}

func ReadJsonUsers() []models.User {
	file, err := os.Open(JsonUsersFile)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var users []models.User
	byteValue, err := io.ReadAll(file)

	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(byteValue, &users)

	return users
}

func WriteUserStore(data []models.User) {
	jsonString, _ := json.Marshal(data)
	os.WriteFile(JsonUsersFile, jsonString, os.ModePerm)
}
