package main

import (
    "FriendsTestTask/Backend/DataBase"
    "FriendsTestTask/Backend/Handlers"
    "log"
    "net/http"
)

func main() {
    																	// Редирект на главную страницу при обращении к корню сайта
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/index.html", http.StatusFound) })

    																	// Обработка запросов к страницам
    http.HandleFunc("/index.html", handlers.HandleIndex)     			// Главная страница
    http.HandleFunc("/logout", handlers.HandleLogout)        			// Выход из аккаунта
    http.HandleFunc("/Registration.html", handlers.HandleRegistration) 	// Страница регистрации
    http.HandleFunc("/main.html", handlers.HandleMain)     				// Основная страница после авторизации
    http.HandleFunc("/style.css", handlers.HandleCSS)        			// Файл стилей

    																	// Запуск сервера на порту 8080
    log.Fatal(http.ListenAndServe(":8080", nil))

    																	// Закрытие соединения с базой данных
    err := database.CloseConnection()
    if err != nil {
        log.Fatal(err)
    }
}