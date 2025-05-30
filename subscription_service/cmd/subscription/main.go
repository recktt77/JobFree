package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"

	"github.com/recktt77/JobFree/subscription_service/config"
	"github.com/recktt77/JobFree/subscription_service/internal/app"
)

func main() {
	// Загрузка переменных окружения из .env
	err := godotenv.Load("C:\\Users\\админ\\Desktop\\HW\\JobFree\\subscription_service\\.env")
	if err != nil {
		log.Println("Error loading .env file:", err)
	}

	log.Println("MONGO_DB_URI:", os.Getenv("MONGO_DB_URI"))

	// Конфигурация
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	// Подключение к NATS
	natsConn, err := nats.Connect(cfg.NATSUrl)
	if err != nil {
		log.Fatalf("failed to connect to NATS: %v", err)
	}
	defer natsConn.Close()

	// Создание контекста
	ctx := context.Background()

	// Инициализация приложения
	application, err := app.New(ctx, cfg, natsConn)
	if err != nil {
		log.Fatalf("failed to setup application: %v", err)
	}

	// Запуск gRPC сервера
	err = application.Run()
	if err != nil {
		log.Fatalf("failed to run application: %v", err)
	}
}
