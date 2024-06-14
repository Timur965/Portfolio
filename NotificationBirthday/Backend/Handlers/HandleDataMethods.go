package handlers

import (
    "FriendsTestTask/Backend/DataBase"
    "errors"
    "io"
    "net/url"
    "strconv"

    "golang.org/x/crypto/bcrypt"
)

														// MeAndUsers - структура, содержащая информацию о текущем пользователе и всех других пользователях.
type MeAndUsers struct {
    MeRecNotif     bool              					// Флаг, указывающий, получает ли текущий пользователь уведомления
    UsersWithoutMe []database.User 						// Список пользователей, не включая текущего пользователя
}

														// Registration - функция регистрации нового пользователя.
func Registration(bodyReader io.ReadCloser) error {
    													// Чтение тела запроса.
    body, err := io.ReadAll(bodyReader)
    if err != nil {
        return err
    }

    													// Разбор URL-кодированных данных.
    data, err := url.ParseQuery(string(body))
    if err != nil {
        return err
    }

    													// Получение данных из запроса.
    username := data.Get("username")
    password := data.Get("password")
    email := data.Get("email")
    firstName := data.Get("first_name")
    lastName := data.Get("last_name")
    birthday := data.Get("birthday")

    													// Хеширование пароля.
    hashPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    													// Вставка нового пользователя в базу данных.
    err = database.InsertUser(username, string(hashPass), email, firstName, lastName, birthday)
    if err != nil {
        return err
    }

    return nil
}

																	// Authorize - функция авторизации существующего пользователя.
func Authorize(bodyReader io.ReadCloser) (database.User, error) {
    																// Чтение тела запроса.
    body, err := io.ReadAll(bodyReader)
    if err != nil {
        return database.User{}, err
    }

    																// Разбор URL-кодированных данных.
    data, err := url.ParseQuery(string(body))
    if err != nil {
        return database.User{}, err
    }

    																// Получение данных из запроса.
    username := data.Get("username")
    pass := data.Get("password")

    																// Получение пользователя из базы данных по имени пользователя.
    user, err := database.GetUser(username)
    if err != nil {
        return database.User{}, err
    }

    																// Проверка пароля.
    if user.Username == username && bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass)) == nil {
        return user, nil
    }

    																// Возврат ошибки, если логин или пароль неверны.
    return database.User{}, errors.New("неправильный логин или пароль")
}

																	// Subscribe - функция подписки на другого пользователя.
func Subscribe(url *url.URL, subscriber string) error {
    																// Получение имени пользователя, на которого подписываются.
    UName := url.Query().Get("Subscribed")
    																// Получение данных о пользователе, который подписывается.
    Usubscriber, err := database.GetUser(subscriber)
    if err != nil {
        return err
    }
    																// Получение данных о пользователе, на которого подписываются.
    Usubscribed, err := database.GetUser(UName)
    if err != nil {
        return err
    }
    																// Вставка информации о подписке в базу данных.
    err = database.InsertSubscription(Usubscriber.Id, Usubscribed.Id, 1)
    if err != nil {
        return err
    }

    return nil
}

																	// DelSubscribed - функция отписки от другого пользователя.
func DelSubscribed(url *url.URL) error {
    																// Получение имени пользователя, от которого отписываются из запроса.
    NameSubscribed := url.Query().Get("DelSubscribed")
    																// Получение данных о пользователе, от которого отписываются.
    Usubscribed, err := database.GetUser(NameSubscribed)
    if err != nil {
        return err
    }
    																// Удаление информации о подписке из базы данных.
    err = database.DeleteSubscribed(Usubscribed.Id)
    if err != nil {
        return err
    }

    return nil
}

																		// DataForMainPage - функция получения данных для главной страницы.
func DataForMainPage(myUser database.User) ([]database.User, error) {
    																	// Получение всех пользователей, кроме текущего.
    userWithoutMe, err := database.GetAllUsersWithoutMe(myUser.Username)
    if err != nil {
        return nil, err
    }

   																		// Получение списка пользователей, на которых подписан текущий пользователь.
    UsersSubsMe, err := database.GetAllUsersSubscribeMe(myUser.Id)
    if err != nil {
        return nil, err
    }

    																	// Цикл для пометки пользователей на которых я подписан в слайсе userWithoutMe
    for i, u := range userWithoutMe {
        if _, ok := UsersSubsMe[u.Username]; ok {
            userWithoutMe[i].Subscribe = true
        }
    }

    return userWithoutMe, nil
}

																	// GetUserFromDB - функция получения данных о пользователе из базы данных.
func GetUserFromDB(username string) (database.User, error) {
    																// Получение пользователя из базы данных по имени пользователя.
    me, err := database.GetUser(username)

    if err != nil {
        return database.User{}, err
    }

    return me, nil
}

																	// UpdateMeNotification - функция обновления настроек уведомлений для текущего пользователя.
func UpdateMeNotification(url *url.URL, me database.User) error {
    																// Получение значения параметра "Notifications" из URL.
    param1 := url.Query().Get("Notifications")
    																// Преобразование значения параметра "Notifications" в булево значение.
    recieveNotification, err := strconv.ParseBool(param1)
    if err != nil {
        return err
    }

    																// Обновление настроек уведомлений, если включено получение уведомлений.
    if recieveNotification {
        															// Получение значения параметра "NotificationDays" из URL.
        param2 := url.Query().Get("NotificationDays")
        															// Преобразование значения параметра "NotificationDays" в целое число.
        notificationTime, err := strconv.ParseInt(param2, 10, 32)
        if err != nil {
            return err
        }

        															// Обновление настроек уведомлений в базе данных.
        err = database.UpdateSubscriptions(int(notificationTime), me.Id)
        if err != nil {
            return err
        }
    }

    																// Обновление значения "RecieveNotifications" в базе данных.
    return database.UpdateReceiveNotifications(me.Username, recieveNotification)
}