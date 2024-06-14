package handlers

import (
    "FriendsTestTask/Backend/DataBase"
    "net/http"
    "time"
)

																	// CreateCookie создает cookie для аутентификации пользователя
func CreateCookie(w http.ResponseWriter, username string) {
    cookie := &http.Cookie{
        Name:     "user_data", 										// Имя cookie
        Value:    username,       									// Значение cookie (имя пользователя)
        Path:     "/",            									// Путь, по которому cookie доступен
        HttpOnly: true,           									// Cookie доступен только для сервера
        Secure:   true,           									// Cookie доступен только по HTTPS
        MaxAge:   int(3600 * 24), 									// Время жизни cookie (24 часа)
    }

    http.SetCookie(w, cookie) 										// Отправляет cookie в браузер
}

																	// CheckAuthorize проверяет, авторизован ли пользователь
func CheckAuthorize(r *http.Request) error {
    cookie, err := r.Cookie("user_data") 							// Получает cookie из запроса
    if err != nil {
        return err 													// Возвращает ошибку, если cookie не найден
    }

    if cookie.Expires.Before(time.Now()) { 							// Проверяет, не истекло ли время жизни cookie
        return err 													// Возвращает ошибку, если время жизни cookie истекло
    }

    _, err = database.GetUser(cookie.Value) 						// Получает пользователя из базы данных
    if err != nil {
        return err 													// Возвращает ошибку, если пользователь не найден в базе данных
    }

    return nil 														// Возвращает nil, если пользователь авторизован
}

																	// GetUserNameFromCookie получает имя пользователя из cookie
func GetUserNameFromCookie(r *http.Request) (string, error) {
    cookie, err := r.Cookie("user_data") 							// Получает cookie из запроса
    if err != nil {
        return "", err 												// Возвращает пустую строку и ошибку, если cookie не найден
    }
    return cookie.Value, nil 										// Возвращает имя пользователя из cookie
}

																	// DeleteCookie удаляет cookie
func DeleteCookie(w http.ResponseWriter, cookieName string) {
    http.SetCookie(w, &http.Cookie{
        Name:     cookieName, 										// Имя cookie, которое нужно удалить
        Value:    "",          										// Пустое значение cookie
        Path:     "/",          									// Путь, по которому cookie доступен
        Expires:  time.Unix(0, 0), 									// Время истечения срока действия cookie
        HttpOnly: true,          									// Cookie доступен только для сервера
    })
}