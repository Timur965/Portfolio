package main

import (
	database "FriendsTestTask/Backend/DataBase"  					// Импортируем пакет "DataBase" из пакета "FriendsTestTask/Backend"
	Send "FriendsTestTask/SendNotification/Send" 					// Импортируем пакет "Send" из пакета "FriendsTestTask/SendNotification"
	"fmt"

	"github.com/robfig/cron/v3" 									// Импортируем пакет "cron" для работы с планировщиком заданий
)

func main() {
    c := cron.New() 												// Создаем новый планировщик заданий
    																// Добавляем задачу в планировщик. Задача будет запускаться каждый день в 00:00
    if _, err := c.AddFunc("0 0 * * *", func() {
        users, err := database.GetUsersForNotification() 			// Получаем список пользователей, которым нужно отправить уведомление
        if err != nil {
            fmt.Println(err) 	 									// Выводим ошибку, если произошла ошибка при получении списка пользователей
        }
        if len(users) != 0 {
            err := Send.SendEmailNotification(users) 				// Отправляем уведомление по электронной почте всем пользователям из списка
			if err != nil {
				fmt.Println(err) 									// Выводим фатальную ошибку, если произошла ошибка при закрытии соединения
			}
        } else {
            fmt.Println("Нет пользователей для отправки")
        }
    }); err != nil {
        fmt.Println("Ошибка добавления задачи в cron:", err) 		// Выводим ошибку, если произошла ошибка при добавлении задачи в планировщик
    }

    err := database.CloseConnection() 								// Закрываем соединение с базой данных
    if err != nil {
        fmt.Println(err) 											// Выводим фатальную ошибку, если произошла ошибка при закрытии соединения
    }

    c.Start() 														// Запускаем планировщик заданий

    select {} 														// Блокируем выполнение программы до тех пор, пока не будет получено событие
}