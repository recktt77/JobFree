package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"github.com/recktt77/JobFree/config"
	"github.com/recktt77/JobFree/internal/app"
)

func main() {
	// Загружаем .env из корня auth_service
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Не удалось загрузить .env файл:", err)
	}

	// Загружаем конфигурацию из переменных
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Ошибка конфигурации: %v", err)
	}

	ctx := context.Background()
	app.Run(ctx, *cfg)
}
