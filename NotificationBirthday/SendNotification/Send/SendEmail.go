package SendNotification

import (
	"FriendsTestTask/Backend/ENV"                                                       // Импорт пакета ENV для получения конфигурации SMTP
	"crypto/tls"                                                                        // Импорт пакета tls для настройки TLS
	"fmt"                                                                               // Импорт пакета fmt для форматирования строк

	"gopkg.in/gomail.v2"                                                                // Импорт пакета gomail для отправки почты
)

                                                                                        // SendEmailNotification отправляет уведомления по электронной почте всем пользователям.
                                                                                        // Принимает карту пользователей, где ключ - это адрес электронной почты, а значение - текст сообщения.
func SendEmailNotification(users map[string]string) error {
    smtpConfig, err := ENV.GetInfoEmail() 												// Получение конфигурации из пакета ENV
    if err != nil {
        return fmt.Errorf("Не удалось получить конфигурацию SMTP: %w", err) 			// Возврат ошибки, если не удалось получить конфигурацию
    }

    if len(smtpConfig.Email) == 0 || len(smtpConfig.Pass) == 0 {                        //Проверка на пустоту майла или пароля
        return fmt.Errorf("Пустой email отправителя или пароль. Проверьте файл env")
    }
    dialer := gomail.NewDialer("smtp.mail.ru", 465, smtpConfig.Email, smtpConfig.Pass) 	// Создание нового диалера с использованием конфигурации
    dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true} 							// Настройка TLS-конфигурации с отключением проверки сертификатов

    message := gomail.NewMessage() 														// Создание нового сообщения
    for email, value := range users { 													// Итерация по всем пользователям в карте
        message.SetHeader("From", smtpConfig.Email) 									// Установка адреса отправителя
        message.SetHeader("To", email)               									// Установка адреса получателя
        message.SetHeader("Subject", "Уведомление о днях рождения") 					// Установка темы сообщения
        message.SetBody("text/plain", value)           									// Установка тела сообщения

        if err := dialer.DialAndSend(message); err != nil { 							// Отправка сообщения по электронной почте
            return fmt.Errorf("Не удалось отправить письмо на  %s: %w", email, err) 	// Возврат ошибки, если не удалось отправить письмо
        }
    }

    return nil 																			// Возврат nil, если все письма отправлены успешно
}
