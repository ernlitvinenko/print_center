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

### Response (Error - 400 Bad Request)

```json
{
  "success": false,
  "error": "phone number is required"
}
```

### Response (Error - 401 Unauthorized)

```json
{
  "success": false,
  "error": "Invalid phone number or password"
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
