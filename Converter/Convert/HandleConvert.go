package Convert

import (
    "Converter/ParseJson"
    "errors"
    "io"
    "math"
    "net/url"
    "strconv"
)

type ResultDataCurrency struct {																		// ResultDataCurrency представляет данные ответа на запрос конвертации валют.
    Success bool																						// Success указывает на успешность операции конвертации.
    Valute1 float64																						// Valute1 - исходное значение валюты.
    Valute2 float64																						// Valute2 - преобразованное значение валюты.
    Curr map[string]ParseJson.Valute																	// Curr - карта валют, полученная из JSON файла.
}

func HandlePOST(rBody io.ReadCloser, curr map[string]ParseJson.Valute) (ResultDataCurrency, error) {	// HandlePOST обрабатывает POST запрос конвертации валют.
    body, err := io.ReadAll(rBody)																		// Читаем тело запроса.
    if err != nil {
        return ResultDataCurrency{}, errors.Join(errors.New("Ошибка чтения тела запроса:"), err)
    }

    data, err := url.ParseQuery(string(body))															// Парсим тело запроса как URL-кодированные данные.
    if err != nil {
        return ResultDataCurrency{}, errors.Join(errors.New("Ошибка парсинга тела запроса:"), err)
    }

    var res ResultDataCurrency																			// Создаем структуру для ответа.
    res.Valute1, err = strconv.ParseFloat(data.Get("ValValute1"), 64)									// Парсим исходное значение валюты.
    if err != nil {
        return ResultDataCurrency{}, errors.Join(errors.New("Ошибка преобразования значения ValValute1 в число:"), err)
    }

																										// Получаем коды валют из запроса.
    Code1 := data.Get("SelectValute1")
    Code2 := data.Get("SelectValute2")

    																									// Проверяем, существуют ли валюты с указанными кодами в карте валют.
    if _, ok := curr[Code1]; !ok {
        return ResultDataCurrency{}, errors.New("Валюта с кодом " + Code1 + " не найдена")
    }
    if _, ok := curr[Code2]; !ok {
        return ResultDataCurrency{}, errors.New("Валюта с кодом " + Code2 + " не найдена")
    }

																										// Выполняем конвертацию валют.
    if Code1 != Code2 {
        if Code1 == "RUB" {																				// Если конвертируем из рублей.
            res.Valute2 = res.Valute1 / curr[Code2].Value
        } else if Code2 == "RUB" {																		// Если конвертируем в рубли.
            res.Valute2 = res.Valute1 * curr[Code1].Value
        } else {																						// Если конвертируем между валютами, не являющимися рублями.
            res.Valute2 = (res.Valute1 * curr[Code1].Value) / (curr[Code2].Value / float64(curr[Code2].Nominal))
        }
    } else {																							// Если коды валют совпадают, просто возвращаем исходное значение.
        res.Valute2 = res.Valute1
    }

    res.Valute2 = math.Round(res.Valute2*100) / 100														// Округляем результат конвертации до сотых.
    res.Success = true

    res.Curr = curr																						// Копируем карту валют в ответ.

    																									// Отмечаем валюты, использованные в конвертации, в карте валют.
    item := res.Curr[Code1]
    item.ActiveCode1 = true
    res.Curr[Code1] = item

    item = res.Curr[Code2]
    item.ActiveCode2 = true
    res.Curr[Code2] = item

    return res, nil																						// Возвращаем результат конвертации.
}