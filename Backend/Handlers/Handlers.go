package handlers

import (
	"html/template"
	"net/http"
	"os"
)

// downloadTemplate - функция для загрузки и обработки шаблона
func downloadTemplate(w http.ResponseWriter, pathTemplate string, data any) {
    																			// Парсинг шаблона из указанного файла
    t := template.Must(template.ParseFiles(pathTemplate))

    																			// Выполнение шаблона с переданными данными и отправка результата в ответ
    err := t.Execute(w, data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
}

															// checkAuthorizedOrRedirect - функция для проверки авторизации пользователя и перенаправления при необходимости
func checkAuthorizedOrRedirect(w http.ResponseWriter, r *http.Request, templatePath string, data interface{}) {
    														// Проверка авторизации пользователя
    if CheckAuthorize(r) != nil {
        													// Если пользователь не авторизован, загружается шаблон и передаются данные
        downloadTemplate(w, templatePath, data)
        return
    }
    														// Перенаправление на главную страницу
    http.Redirect(w, r, "/main.html", http.StatusFound)
}

																// HandleIndex - обработчик для авторизации
func HandleIndex(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    															// Обработка POST-запроса для авторизации
    if r.Method == "POST" {
        														// Аутентификация пользователя с помощью данных из тела запроса
        user, err := Authorize(r.Body)

        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        } else {
            													// Создание куки с именем пользователя
            CreateCookie(w, user.Username)
            													// Перенаправление на главную страницу
            http.Redirect(w, r, "/main.html", http.StatusFound)
            return
        }
    } else {
        														// Вызов функции проверки авторизации и загрузки шаблона
        checkAuthorizedOrRedirect(w, r, "index.html", nil)
    }
}

															// HandleCSS - обработчик для загрузки CSS-файла
func HandleCSS(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    														// Чтение содержимого CSS-файла
    data, err := os.ReadFile("style.css")
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    														// Установка заголовка для типа контента
    w.Header().Set("Content-Type", "text/css")

    														// Отправка содержимого файла в ответ
    w.Write(data)
}

															// HandleRegistration - обработчик для страницы регистрации
func HandleRegistration(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    														// Обработка POST-запроса для регистрации
    if r.Method == "POST" {
        													// Регистрация пользователя с помощью данных из тела запроса
        err := Registration(r.Body)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        } else {
            												// Перенаправление на страницу входа если регестрация успешна
            http.Redirect(w, r, "/index.html", http.StatusFound)
            return
        }
    } else {
        													// Вызов функции проверки авторизации и загрузки шаблона
        checkAuthorizedOrRedirect(w, r, "Registration.html", nil)
    }
}

															// HandleMain - обработчик для главной страницы после авторизации
func HandleMain(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    														// Проверка авторизации пользователя
    if CheckAuthorize(r) != nil {
        													// Перенаправление на страницу входа
        http.Redirect(w, r, "/index.html", http.StatusFound)
        return
    } else {
        													// Создание структуры данных для шаблона
        var info MeAndUsers
        													// Получение имени пользователя из куки
        usernameFromCookie, err := GetUserNameFromCookie(r)
        if err != nil {
            http.Error(w, err.Error(), http.StatusUnauthorized)
            return
        }

        													// Получение данных о пользователе из базы данных
        myUser, err := GetUserFromDB(usernameFromCookie)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        													// Заполнение данных о пользователе в структуру
        info.MeRecNotif = myUser.RecNotif

        													// Получение данных о других пользователях для отображения на странице
        info.UsersWithoutMe, err = DataForMainPage(myUser)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        													// Обработка запросов GET для подписки на пользователей
        if r.Method == "GET" && r.URL.Query().Has("Subscribed") {
            												// Подписка на пользователя
            err = Subscribe(r.URL, usernameFromCookie)
            if err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            } else {
                											// Перенаправление на главную страницу
                http.Redirect(w, r, "/main.html", http.StatusFound)
                return
            }
        }

        													// Обработка запросов GET для отписки от пользователей
        if r.Method == "GET" && r.URL.Query().Has("DelSubscribed") {
            												// Отписка от пользователя
            err = DelSubscribed(r.URL)
            if err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            } else {
                											// Перенаправление на главную страницу
                http.Redirect(w, r, "/main.html", http.StatusFound)
                return
            }
        }

        													// Обработка запросов GET для изменения настроек уведомлений
        if r.Method == "GET" && r.URL.Query().Has("Notifications") {
            												// Обновление настроек уведомлений пользователя
            err = UpdateMeNotification(r.URL, myUser)
            if err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            } else {
                											// Перенаправление на главную страницу
                http.Redirect(w, r, "/main.html", http.StatusFound)
                return
            }
        }

        													// Загрузка и обработка шаблона с переданными данными
        downloadTemplate(w, "main.html", info)
    }
}

															// HandleLogout - обработчик для выхода из системы
func HandleLogout(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    														// Удаление куки с именем пользователя
    DeleteCookie(w, "user_data")
    														// Перенаправление на страницу входа
    http.Redirect(w, r, "/index.html", http.StatusFound)
}