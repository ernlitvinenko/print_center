package main

import (
	"backend/core/services"
	"fmt"
	"os"
)

// Утилита для генерации bcrypt хеша пароля
// Использование: go run cmd/hash_password/main.go <ваш_пароль>
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: go run cmd/hash_password/main.go <пароль>")
		os.Exit(1)
	}

	password := os.Args[1]
	hashedPassword, err := services.HashPassword(password)
	if err != nil {
		fmt.Printf("Ошибка при хешировании пароля: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Пароль: %s\n", password)
	fmt.Printf("Bcrypt хеш: %s\n", hashedPassword)
	fmt.Println("\nИспользуйте этот хеш для вставки в поле password таблицы profile")
}
