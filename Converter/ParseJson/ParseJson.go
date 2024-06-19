package ParseJson

import (
	"encoding/json"
	"errors"
	"time"
)

															// data структура для представления JSON-данных
type data struct {
	Date   time.Time         `json:"Date"`
	Valute map[string]Valute `json:"Valute"`
}

															// Valute структура для представления данных о валюте
type Valute struct {
	Nominal     int     `json:"Nominal"`
	Name        string  `json:"Name"`
	Value       float64 `json:"Value"`
	ActiveCode1 bool
	ActiveCode2 bool
}

															// Parse функция для парсинга JSON-данных
func Parse(jsonData []byte) (map[string]Valute, error) {
	var data data
	err := json.Unmarshal([]byte(jsonData), &data)

	if err != nil {
		return nil, errors.Join(errors.New("Ошибка десериализации JSON"), err)
	}

	return data.Valute, nil
}
