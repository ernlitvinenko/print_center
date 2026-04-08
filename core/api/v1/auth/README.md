# Авторизация API

## POST /api/v1/auth/login

Вход пользователя по номеру телефона и паролю.

### Request Body

```json
{
  "phone": "+375 (29) 123-45-67",
  "password": "your_password"
}
```

### Response (Success - 200 OK)

```json
{
  "success": true,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "first_name": "Иван",
    "last_name": "Иванов",
    "email": "ivan@example.com",
    "phone": "375291234567"
  }
}
```

## GET /api/v1/auth/me

Возвращает информацию о текущем пользователе (требует JWT токен).

### Headers

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Response (Success - 200 OK)

```json
{
  "success": true,
  "user": {
    "id": 1,
    "first_name": "Иван",
    "last_name": "Иванов",
    "phone": "375291234567"
  }
}
```

### Response (Error - 401 Unauthorized)

```json
{
  "error": "invalid or expired token"
}
```

## Регистрация нового пользователя

Для генерации bcrypt хеша пароля используйте утилиту:

```bash
go run cmd/hash_password/main.go my_secure_password
```

Затем вставьте полученный хеш в базу данных:

```sql
INSERT INTO profile (first_name, last_name, father_name, email, phone_dgt, password)
VALUES ('Иван', 'Иванов', 'Иванович', 'ivan@example.com', 375291234567, '$2a$10$...');
```

## Валидация

- Номер телефона: обязателен, 10-15 цифр
- Пароль: обязателен, минимум 6 символов

## Безопасность

- Пароли хранятся в базе данных в захешированном виде (bcrypt)
- Bcrypt автоматически использует salt для защиты от rainbow table атак
- Стоимость хеширования: `bcrypt.DefaultCost` (10)

## Использование JWT токенов

### 1. Получение токена

Отправьте POST запрос на `/api/v1/auth/login`:

```bash
curl -X POST http://localhost:8000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "+375 (29) 123-45-67",
    "password": "my_secure_password"
  }'
```

В ответе получите токен в поле `token`.

### 2. Использование токена

Добавьте токен в заголовок `Authorization`:

```bash
curl http://localhost:8000/api/v1/auth/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 3. Конфигурация JWT

В `.env` файле можно настроить:

```env
# Секретный ключ для подписи токенов (обязательно измените в production!)
JWT_SECRET=your-super-secret-key-change-this-in-production
```

### 4. Время жизни токена

По умолчанию токен действует **24 часа**. Для изменения отредактируйте `core/services/jwt_service.go`:

```go
func NewJWTConfig() *JWTConfig {
    return &JWTConfig{
        SecretKey:     getEnv("JWT_SECRET", "..."),
        TokenDuration: time.Hour * 24, // Измените здесь
    }
}
```

## Структура

```
core/
├── api/v1/auth/          # HTTP обработчики
│   └── login.go          # Login endpoint
├── models/               # Общие модели
│   ├── request.go        # LoginRequest
│   └── response.go       # LoginResponse, UserInfo
└── services/             # Бизнес-логика
    └── auth_service.go   # AuthService с bcrypt
```
