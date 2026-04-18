package main

import (
	"async-api-task-manager/internal/storage/postgres"
	"fmt"
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

	// graceful shutdown
}

func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
