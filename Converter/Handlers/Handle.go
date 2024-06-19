package Handlers

import (
	"Converter/Convert"
	"Converter/Currency"
	"errors"
	"html/template"
	"net/http"
)

func HandleIndex(w http.ResponseWriter, r *http.Request) { 						// HandleIndex функция обработки запроса к корневому пути
	t, err := template.ParseFiles("index.html") 								// Парсим шаблон index.html
	if err != nil {																// Если возникла ошибка при парсинге шаблона, выводим ее пользователю
		http.Error(w, errors.Join(errors.New("Ошибка при парсинге шаблона:"), err).Error(), http.StatusInternalServerError)
		return 																	// Прерываем обработку запроса
	}
	resCurr := Convert.ResultDataCurrency{}										// Инициализируем структуру данных для отображения валют
	resCurr.Curr, err = Currency.PrintCurrencyData()							// Получаем актуальные данные о валютах

	if err != nil {																// Если возникла ошибка при получении данных о валютах, выводим ее пользователю
		http.Error(w, errors.Join(errors.New("Ошибка при получении данных о валютах:"), err).Error(), http.StatusInternalServerError)
		return 																	// Прерываем обработку запроса
	}

	if r.Method == http.MethodPost {											// Обрабатываем POST-запрос
		resCurr, err = Convert.HandlePOST(r.Body, resCurr.Curr)					// Вызываем функцию HandlePOST для обработки данных из POST-запроса
		if err != nil {															// Если возникла ошибка при обработке POST-запроса, выводим ее пользователю
			http.Error(w, errors.Join(errors.New("Ошибка при обработке POST-запроса:"), err).Error(), http.StatusBadRequest)
			return 																// Прерываем обработку запроса
		}
	}

	err = t.Execute(w, resCurr)													// Выполняем шаблон, передавая данные о валютах
	if err != nil {																// Если возникла ошибка при выполнении шаблона, выводим ее пользователю
		http.Error(w, errors.Join(errors.New("Ошибка при отображении шаблона:"), err).Error(), http.StatusInternalServerError)
		return 																	// Прерываем обработку запроса
	}
}
