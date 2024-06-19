FROM golang:latest 

WORKDIR /app

# Сначала копируем файлы зависимостей
COPY . .

# Сборка приложения
RUN go build -o bin .

# Задаем точку входа
ENTRYPOINT ["/app/bin"] 