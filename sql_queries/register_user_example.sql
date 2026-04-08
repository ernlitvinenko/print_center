-- Пример регистрации пользователя с bcrypt хешем пароля
-- Замените <bcrypt_hash> на реальный хеш, сгенерированный через: go run cmd/hash_password/main.go <пароль>

INSERT INTO profile (first_name, last_name, father_name, email, phone_dgt, password)
VALUES (
    'Иван',
    'Иванов',
    'Иванович',
    'ivan@example.com',
    375291234567,
    '<bcrypt_hash>'  -- Вставьте сюда bcrypt хеш
);
