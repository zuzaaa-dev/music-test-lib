# Используем официальный образ Golang
FROM golang:1.23-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы модуля и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем файл .env в рабочую директорию
COPY .env .env

# Копируем весь код в рабочую директорию
COPY . .

# Копируем миграции в рабочую директорию
COPY migrations/ ./migrations/

# Компилируем приложение
RUN go build -o music ./cmd/music/main.go

# Используем минимальный образ для запуска
FROM alpine:latest

# Копируем скомпилированное приложение из предыдущего этапа
WORKDIR /root/
COPY --from=builder /app/music .
# Копируем .env файл в конечный образ
COPY --from=builder /app/.env .
# Копируем миграции в конечный образ
COPY --from=builder /app/migrations ./migrations

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./music"]
