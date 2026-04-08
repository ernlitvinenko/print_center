# Настройка проекта Print Center Backend

## Требования

- Go 1.25+
- PostgreSQL 14+

## Настройка базы данных

### 1. Установка PostgreSQL

Если PostgreSQL не установлен:

```bash
# macOS
brew install postgresql@16
brew services start postgresql@16

# Или через Postgres.app
# https://postgresapp.com/
```

### 2. Создание базы данных

```bash
# Подключаемся к PostgreSQL
psql -U postgres

# Создаём базу данных и пользователя
CREATE DATABASE print_center;
\c print_center

# Выполняем схему из файла
\i sql_queries/schema.sql

# Или из командной строки:
psql -U postgres -d print_center -f sql_queries/schema.sql
```

### 3. Создание тестового пользователя

```bash
# Генерируем bcrypt хеш пароля
go run cmd/hash_password/main.go my_secure_password

# Вставляем пользователя (замените <bcrypt_hash> на реальный хеш)
psql -U postgres -d print_center -c "
INSERT INTO profile (first_name, last_name, father_name, email, phone_dgt, password)
VALUES (
    'Иван',
    'Иванов',
    'Иванович',
    'ivan@example.com',
    375291234567,
    '\$2a\$10\$...'  -- вставьте сюда реальный хеш
);
"
```

### 4. Настройка переменных окружения

```bash
# Копируем пример файла с переменными
cp .env.example .env

# Редактируем .env под свою конфигурацию
nano .env
```

Содержимое `.env`:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=print_center

APP_PORT=8000
LOG_LEVEL=info
```

## Запуск проекта

```bash
# Устанавливаем зависимости
go mod tidy

# Запускаем сервер
go run main.go
```

Сервер запустится на `http://localhost:8000`

## API Endpoints

### Health Check

```bash
curl http://localhost:8000/
```

Response:
```json
{
  "status": "ok",
  "message": "Print Center API is running"
}
```

### Login

```bash
curl -X POST http://localhost:8000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "+375 (29) 123-45-67",
    "password": "my_secure_password"
  }'
```

Response (200 OK):
```json
{
  "success": true,
  "token": "",
  "user": {
    "id": 1,
    "first_name": "Иван",
    "last_name": "Иванов",
    "email": "ivan@example.com",
    "phone": "375291234567"
  }
}
```

Response (401 Unauthorized):
```json
{
  "success": false,
  "error": "Invalid phone number or password"
}
```

## Структура проекта

```
backend/
├── cmd/
│   └── hash_password/        # Утилита для генерации bcrypt хеша
├── config/
│   └── config.go             # Конфигурация и подключение к БД
├── core/
│   ├── api/v1/auth/          # HTTP обработчики авторизации
│   │   └── login.go
│   │       └── README.md
│   ├── models/               # Общие модели запросов/ответов
│   │   ├── request.go
│   │   └── response.go
│   ├── services/             # Бизнес-логика (сервисы)
│   │   └── auth_service.go   # AuthService с bcrypt
│   └── repositories/         # Сгенерировано sqlc
│       ├── db.go
│       ├── models.go
│       └── query.sql.go
├── sql_queries/
│   ├── schema.sql            # Схема базы данных
│   ├── query.sql             # SQL запросы
│   └── register_user_example.sql
├── .env.example              # Пример переменных окружения
├── .gitignore
├── go.mod
├── go.sum
├── main.go
└── sqlc.yaml
```

## Разработка

### Перегенерация sqlc

После изменения `sql_queries/query.sql` или `schema.sql`:

```bash
sqlc generate
```

### Зависимости

```bash
# Добавить зависимости
go get <package>

# Обновить все зависимости
go get -u ./...

# Очистить неиспользуемые зависимости
go mod tidy
```

## Troubleshooting

### Ошибка подключения к БД

```
failed to connect to database: failed to connect to `host=localhost`
```

**Решение:**
1. Убедитесь, что PostgreSQL запущен: `brew services list | grep postgresql`
2. Проверьте переменные окружения в `.env`
3. Проверьте, что база данных существует: `psql -U postgres -l | grep print_center`

### Ошибка аутентификации

```
invalid password
```

**Решение:**
1. Убедитесь, что пароль в БД хранится в bcrypt формате (начинается с `$2a$`, `$2b$` или `$2x$`)
2. Перегенерируйте хеш: `go run cmd/hash_password/main.go <пароль>`
3. Обновите пароль в базе данных
