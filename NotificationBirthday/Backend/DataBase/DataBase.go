package database

import (
	"FriendsTestTask/Backend/ENV" // Импорт пакета ENV для получения настроек БД
	"database/sql"                // Импорт пакета sql для работы с БД
	"fmt"                         // Импорт пакета fmt для форматирования строк
	"time"                        // Импорт пакета time для работы с датами и временем

	_ "github.com/lib/pq" // Импорт драйвера PostgreSQL
)

var globdb *sql.DB 							// Глобальная переменная для хранения соединения с БД

											// Структура User для представления пользователя в БД
type User struct {
    Id         int       					// Идентификатор пользователя
    Username   string    					// Имя пользователя
    Password   string    					// Пароль пользователя
    Email      string    					// Электронная почта пользователя
    First_name string    					// Имя пользователя
    Last_name  string    					// Фамилия пользователя
    Birthday   time.Time 					// День рождения пользователя
    RecNotif   bool      					// Флаг уведомлений
    Subscribe  bool      					// Флаг подписки
}

											// Функция для установления соединения с БД
func connect() error {
    if globdb == nil { 						// Проверка на наличие соединения
        var err error

        resENV, err := ENV.GetInfoDB() 		// Получение настроек БД из пакета ENV
        if err != nil {
            return err
        }
                                                            // Формирование строки подключения
        connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", resENV.Host, resENV.Port, resENV.User, resENV.Pass, resENV.DbName)
        globdb, err = sql.Open("postgres", connStr) 		// Открытие соединения с БД

        if err != nil {
            return err
        }

        err = globdb.Ping() 								// Проверка соединения
        if err != nil {
            return err
        }
    }
    return nil
}

									// Функция для закрытия соединения с БД
func CloseConnection() error {
    if globdb != nil { 				// Проверка на наличие соединения
        return globdb.Close() 		// Закрытие соединения
    }
    return nil
}

																			// Функция для выполнения запроса к БД
func executeQuery(query string, args ...interface{}) (*sql.Rows, error) {
    err := connect() 														// Установка соединения с БД
    if err != nil {
        return nil, err
    }

    stmt, err := globdb.Prepare(query) 										// Подготовка запроса
    if err != nil {
        return nil, err
    }
    defer stmt.Close() 														// Закрытие подготовленного запроса

    rows, err := stmt.Query(args...) 										// Выполнение запроса
    if err != nil {
        return nil, err
    }

    return rows, nil
}

																																// Функция для добавления нового пользователя в БД
func InsertUser(username, password, email, first_name, last_name, birthday string) error {
    query := "INSERT INTO users (username, password, email, first_name, last_name, birthday) VALUES ($1, $2, $3, $4, $5, $6)" 	// Запрос на добавление пользователя
    _, err := executeQuery(query, username, password, email, first_name, last_name, birthday) 									// Выполнение запроса
    if err != nil {
        return err
    }

    return nil
}

																															// Функция для получения пользователя по имени пользователя
func GetUser(username string) (User, error) {
    query := "SELECT * FROM users WHERE username = $1" 																		// Запрос на получение пользователя
    rows, err := executeQuery(query, username) 																				// Выполнение запроса
    if err != nil {
        return User{}, err
    }
    defer rows.Close() 																										// Закрытие результата запроса

    u := User{} 																											// Создание структуры User для хранения данных пользователя
    if rows.Next() { 																										// Проверка наличия результата
        err = rows.Scan(&u.Id, &u.Username, &u.Password, &u.Email, &u.First_name, &u.Last_name, &u.Birthday, &u.RecNotif) 	// Считывание данных пользователя
        if err != nil {
            return User{}, err
        }
    } else {
        return User{}, fmt.Errorf("пользователь с именем '%s' не найден", username) 										// Обработка случая, если пользователь не найден
    }

    return u, nil
}

																															// Функция для получения всех пользователей, кроме указанного
func GetAllUsersWithoutMe(username string) ([]User, error) {
    query := "SELECT * FROM users WHERE username != $1" 																	// Запрос на получение всех пользователей, кроме указанного
    rows, err := executeQuery(query, username) 																				// Выполнение запроса
    if err != nil {
        return nil, err
    }
    defer rows.Close() 																										// Закрытие результата запроса

    users := []User{} 																										// Создание среза для хранения данных пользователей
    var u User
    for rows.Next() { 																										// Проверка наличия результата
        err = rows.Scan(&u.Id, &u.Username, &u.Password, &u.Email, &u.First_name, &u.Last_name, &u.Birthday, &u.RecNotif) 	// Считывание данных пользователя
        if err != nil {
            return nil, err
        }
        users = append(users, u) 																							// Добавление пользователя в срез
    }

    return users, nil
}

																															// Функция для получения всех пользователей, на которых подписан пользователь
func GetAllUsersSubscribeMe(UserID int) (map[string]User, error) {
    query := "SELECT u.username FROM users u JOIN subscriptions s ON u.id = s.subscribed_id WHERE s.subscriber_id = $1" 	// Запрос на получение всех пользователей, на которых подписан пользователь
    rows, err := executeQuery(query, UserID) 																				// Выполнение запроса
    if err != nil {
        return nil, err
    }
    defer rows.Close() 																										// Закрытие результата запроса

    users := make(map[string]User) 																							// Создание словаря для хранения данных пользователей
    var u User
    for rows.Next() { 																										// Проверка наличия результата
        err = rows.Scan(&u.Username) 																						// Считывание имени пользователя
        if err != nil {
            return nil, err
        }
        users[u.Username] = u 																								// Добавление пользователя в словарь
    }

    return users, nil
}

																												// Функция для добавления подписки на пользователя
func InsertSubscription(subscriberId, subscribedId, notificationTime int) error {
    query := "INSERT INTO subscriptions (subscriber_id, subscribed_id, notification_time) VALUES ($1, $2, $3)" 	// Запрос на добавление подписки
    _, err := executeQuery(query, subscriberId, subscribedId, notificationTime) 								// Выполнение запроса
    if err != nil {
        return err
    }

    return nil
}

																	// Функция для удаления подписки на пользователя
func DeleteSubscribed(subscribedId int) error {
    query := "DELETE FROM subscriptions WHERE subscribed_id = $1" 	// Запрос на удаление подписки
    _, err := executeQuery(query, subscribedId) 					// Выполнение запроса
    if err != nil {
        return err
    }

    return nil
}

																					// Функция для обновления флага уведомлений пользователя
func UpdateReceiveNotifications(username string, receiveNotifications bool) error {
    query := "UPDATE users SET receive_notification = $1 WHERE username = $2" 		// Запрос на обновление флага уведомлений
    _, err := executeQuery(query, receiveNotifications, username) 					// Выполнение запроса
    if err != nil {
        return err
    }

    return nil
}

																							// Функция для обновления времени уведомлений для подписчиков
func UpdateSubscriptions(notification_time int, subscribed_id int) error {
    query := "UPDATE subscriptions SET notification_time = $1 WHERE subscriber_id = $2" 	// Запрос на обновление времени уведомлений
    _, err := executeQuery(query, notification_time, subscribed_id) 						// Выполнение запроса
    if err != nil {
        return err
    }

    return nil
}

																																	// Функция для получения пользователей, которым нужно отправить уведомления
func GetUsersForNotification() (map[string]string, error) {
    query := `SELECT 
                                u1.Email AS subscriber_first_name,
                                u2.First_name,
                                u2.Last_name,
                                u2.Birthday
                             FROM users AS u1
                                JOIN subscriptions AS s
                                ON u1.id = s.subscriber_id
                                JOIN users AS u2
                                ON s.subscribed_id = u2.id
                             WHERE
                                u1.id != u2.id 
                                AND u1.receive_notification = true 
                                AND u2.receive_notification = true
                                AND DATE(u2.birthday) = DATE(CURRENT_DATE + INTERVAL '1 day' * s.notification_time);` 				// Запрос на получение пользователей, которым нужно отправить уведомления
    rows, err := executeQuery(query) 																								// Выполнение запроса
    if err != nil {
        return nil, err
    }
    defer rows.Close() 																												// Закрытие результата запроса

    users := make(map[string]string) 																								// Создание словаря для хранения данных пользователей
    var u User
    for rows.Next() { // Проверка наличия результата
        err = rows.Scan(&u.Email, &u.First_name, &u.Last_name, &u.Birthday) 														// Считывание данных пользователя
        if err != nil {
            return nil, err
        }
        users[u.Email] += fmt.Sprintf("%s День рождение у %s %s\n", u.Birthday.Format("02-01-2006"), u.First_name, u.Last_name) 	// Формирование текста уведомления
    }

    return users, nil
}