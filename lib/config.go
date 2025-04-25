package lib

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	JsonStoreFile string
	JsonUsersFile string
	PortNo        string
	JwtKey        []byte
)

func LoadConfig(envFile string) {
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Error loading env file: %v", err)
	}

	JsonStoreFile = os.Getenv("json_filename")
	if JsonStoreFile == "" {
		log.Fatal("json_filename is not set in the environment")
	}

	JsonUsersFile = os.Getenv("json_users_filename")
	if JsonUsersFile == "" {
		log.Fatal("json_users_filename is not set in the environment")
	}

	PortNo = os.Getenv("port_number")
	if PortNo == "" {
		log.Fatal("port_number is not set in the environment")
	}

	JwtKey = []byte(os.Getenv("jwt_key"))
	if len(JwtKey) == 0 {
		log.Fatal("jwt_key is not set in the environment")
	}
}
