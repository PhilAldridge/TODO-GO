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
	BaseUrl       string
	SqlPortNo     string
	SqlUser       string
	SqlPassword   string
	SqlDbName     string
	SqlHost       string
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

	BaseUrl = os.Getenv("base_url")
	if len(BaseUrl) == 0 {
		log.Fatal("base_url is not set in the environment")
	}

	SqlDbName = os.Getenv("sql_db_name")
	if len(SqlDbName) == 0 {
		log.Fatal("sql_db_name is not set in the environment")
	}

	SqlPassword = os.Getenv("sql_password")
	if len(SqlPassword) == 0 {
		log.Fatal("sql_password is not set in the environment")
	}

	SqlPortNo = os.Getenv("sql_port_number")
	if len(SqlPortNo) == 0 {
		log.Fatal("sql_port_number is not set in the environment")
	}

	SqlUser = os.Getenv("sql_username")
	if len(SqlUser) == 0 {
		log.Fatal("sql_username is not set in the environment")
	}

	SqlHost = os.Getenv("sql_host")
	if len(SqlHost) == 0 {
		log.Fatal("sql_host is not set in the environment")
	}
}
