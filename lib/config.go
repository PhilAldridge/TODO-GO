package lib


import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var JsonFile string

func LoadConfig(envFile string) {
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Error loading env file: %v", err)
	}

	JsonFile = os.Getenv("json_filename")
	if JsonFile == "" {
		log.Fatal("json_filename is not set in the environment")
	}
}