package main

import (
	"Converter/Handlers"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", Handlers.HandleIndex)					// Регистрируем обработчик запросов для корневого маршрута "/"
	err := (http.ListenAndServe(":8080", nil))					// Запускаем HTTP-сервер на порту 8080
	if err != nil{												// Проверяем наличие ошибок при запуске сервера
		fmt.Errorf("Ошибка запуска сервера: %s", err.Error())	// Если возникла ошибка, выводим сообщение об ошибке
	}
}
