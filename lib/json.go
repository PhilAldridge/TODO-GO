package lib

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/PhilAldridge/TODO-GO/models"
)

func ReadJson() []models.Todo {
	file, err:= os.Open("todoStore.json")

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
	// file, err:= os.Open("todoStore.json")

	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer file.Close()


	// encoder:= json.NewEncoder(file)
	// encoder.Encode(data)
	jsonString,_:= json.Marshal(data)
	os.WriteFile("todoStore.json",jsonString,os.ModePerm) 
}