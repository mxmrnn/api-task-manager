package main

import (
	"async-api-task-manager/internal/service"
	"async-api-task-manager/internal/storage/postgres"
	repository2 "async-api-task-manager/internal/storage/postgres/repository"
	"async-api-task-manager/internal/transport/http"
	"async-api-task-manager/internal/transport/http/handler"

	"fmt"
	"log"
	"os"
)

func main() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "task_manager"),
		getEnv("DB_PORT", "5432"),
	)

	database, err := postgres.NewPostgres(dsn)

	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	log.Println("database connected")

	taskRepo := repository2.NewTaskRepository(database)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := handler.NewTaskHandler(taskService)

	r := http.SetupRouter(taskHandler)
	log.Println("router initialized")

	log.Println("server started on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("run server: %v", err)
	}
	// graceful shutdown
}

func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
