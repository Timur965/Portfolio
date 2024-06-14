package ENV

import (
	"os"

	"github.com/joho/godotenv"
)

type infoForDB struct {
	User   string
	Pass   string
	DbName string
	Host   string
	Port   string
}

type infoForEmail struct {
	Email string
	Pass  string
}

func GetInfoDB() (infoForDB, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return infoForDB{}, err
	}
	return infoForDB{os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT")}, nil
}

func GetInfoEmail() (infoForEmail, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return infoForEmail{}, err
	}
	return infoForEmail{os.Getenv("EMAIL"), os.Getenv("PASS_EMAIL")}, nil
}
