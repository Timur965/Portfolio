package ENV

import (
    "os"

    "github.com/joho/godotenv"
)

												// infoForDB представляет информацию для подключения к базе данных
type infoForDB struct {
    User   string
    Pass   string
    DbName string
    Host   string
    Port   string
}

												// infoForEmail представляет информацию для отправки электронных писем
type infoForEmail struct {
    Email string
    Pass  string
}

												// GetInfoDB извлекает информацию о подключении к базе данных из файла .env
func GetInfoDB() (infoForDB, error) {
    											// Загружаем переменные окружения из файла .env
    err := godotenv.Load(".env")
    if err != nil {
        return infoForDB{}, err
    }
    											// Возвращаем информацию о подключении к базе данных
    return infoForDB{os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT")}, nil
}

												// GetInfoEmail извлекает информацию об отправке электронных писем из файла .env
func GetInfoEmail() (infoForEmail, error) {
    											// Загружаем переменные окружения из файла .env
    err := godotenv.Load(".env")
    if err != nil {
        return infoForEmail{}, err
    }
    											// Возвращаем информацию об отправке электронных писем
    return infoForEmail{os.Getenv("EMAIL"), os.Getenv("PASS_EMAIL")}, nil
}