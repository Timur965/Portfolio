package Currency

import (
    "Converter/ParseJson"
    "errors"
    "io"
    "net/http"
)


func PrintCurrencyData() (map[string]ParseJson.Valute, error) {										// PrintCurrencyData - функция для получения и вывода данных о валютах из API ЦБ РФ
    jsonData, err := getCurrencies("https://www.cbr-xml-daily.ru/daily_json.js")					// Получаем JSON данные с API ЦБ РФ
    if err != nil {																					// Если произошла ошибка при получении данных, возвращаем ошибку
        return nil, errors.Join(errors.New("Не удалось получить данные о валютах:"), err)
    }

    curr, err := ParseJson.Parse(jsonData)															// Парсим полученные JSON данные
    if err != nil {																					// Если произошла ошибка при парсинге данных, возвращаем ошибку
        return nil, errors.Join(errors.New("Не удалось распарсить данные о валютах:"), err)
    }

    curr["RUB"] = ParseJson.Valute{0, "Российский рубль", 0, false, false}                          // Добавили рубль
    return curr, nil																				// Возвращаем полученные данные о валютах
}


func getCurrencies(url string) ([]byte, error) {													// getCurrencies - функция для отправки HTTP запроса к API ЦБ РФ и получения JSON данных
    resp, err := http.Get(url)																		// Отправляем HTTP GET запрос к указанному URL

    if err != nil {																					// Если произошла ошибка при отправке запроса, возвращаем ошибку
        return []byte{}, errors.Join(errors.New("Не удалось выполнить HTTP запрос к API:"), err)
    }

    defer resp.Body.Close()																			// Закрываем тело ответа после завершения работы с ним
    result, err := io.ReadAll(resp.Body)															// Читаем тело ответа в байтовый массив

    if err != nil {																					// Если произошла ошибка при чтении тела ответа, возвращаем ошибку
        return []byte{}, errors.Join(errors.New("Не удалось прочитать тело ответа:"), err)
    }

    return result, nil																				// Возвращаем полученные JSON данные
}