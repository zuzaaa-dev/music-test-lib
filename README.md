Пример файла .env:
```
# Переменные приложения
ENV=local

# HTTP server конфигурация
HTTP_SERVER_ADDRESS=0.0.0.0:8080

# Database configuration
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=postgres
DB_FILE_MIGRATIONS=file://migrations/

# Внешний API
API_MUSIC_INFO_URL=https://localhost:8080/info

# PostgreSQL БД конфигурация
POSTGRES_USER=postgres
POSTGRES_PASSWORD=secret
POSTGRES_DB=postgres
```
