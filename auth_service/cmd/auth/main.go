package main

import (
    "github.com/recktt77/JobFree/config"
    "github.com/recktt77/JobFree/internal/app"
    "context"
    "github.com/joho/godotenv"
    "log"
)

func main() {
    // Загружаем .env из корня auth_service
    if err := godotenv.Load("C:\\Users\\админ\\Desktop\\HW\\JobFree\\auth_service\\.env"); err != nil {
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
