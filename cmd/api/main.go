package main

import (
	"async-api-task-manager/internal/client"
	"async-api-task-manager/internal/service"
	"async-api-task-manager/internal/storage/postgres"
	repository2 "async-api-task-manager/internal/storage/postgres/repository"
	"async-api-task-manager/internal/transport/http"
	"async-api-task-manager/internal/transport/http/handler"
	"async-api-task-manager/internal/transport/rpc"
	"async-api-task-manager/internal/worker"
	"context"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w := worker.NewWorker(10)
	go w.Start(ctx)

	log.Println("worker launched")

	userRepo := repository2.NewUserRepository(database)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	rpcClient, err := client.NewRPCClient("amqp://guest:guest@rabbitmq:5672")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rpcClient.Close()

	authClient := client.NewAuthClient(rpcClient)

	taskRepo := repository2.NewTaskRepository(database)
	taskService := service.NewTaskService(taskRepo, authClient, w)
	taskHandler := handler.NewTaskHandler(taskService)

	rpcServer, err := rpc.NewServer("amqp://guest:guest@rabbitmq:5672")
	if err != nil {
		log.Fatalf("failed to create rpc server: %v", err)
	}
	defer rpcServer.Close()

	rpcHandler := rpc.NewHandler(taskService)

	err = rpcServer.Start("task-service.rpc", rpcHandler.GetUserTaskStatsHandler)
	if err != nil {
		log.Fatalf("failde to start rpc server: %v", err)
	}

	r := http.SetupRouter(taskHandler, userHandler)
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
