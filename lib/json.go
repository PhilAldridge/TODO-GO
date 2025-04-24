package lib

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/PhilAldridge/TODO-GO/models"
)

func ReadJson() []models.Todo {
	file, err:= os.Open(JsonFile)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var todos []models.Todo
	byteValue,err:= io.ReadAll(file)

	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(byteValue, &todos)

	return todos
}

func WriteJson(data []models.Todo) {
	jsonString,_:= json.Marshal(data)
	os.WriteFile(JsonFile,jsonString,os.ModePerm) 
}