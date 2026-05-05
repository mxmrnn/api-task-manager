//go:build integration

package repository

import (
	"async-api-task-manager/internal/model"
	"context"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"testing"
	"time"

	pgDriver "gorm.io/driver/postgres"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	TestDB      *gorm.DB
	TestRepo    *TaskRepository
	pgContainer *postgres.PostgresContainer
)

func SetupSuite(t *testing.T) {
	if t != nil {
		t.Helper()
	}
	ctx := context.Background()

	var err error
	pgContainer, err = postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("task_manager_test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		fatal(t, "failed to start postgres container: %v", err)
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		fatal(t, "failed to get connection string: %v", err)
	}

	TestDB, err = gorm.Open(pgDriver.Open(connStr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		fatal(t, "failed to connect to database: %v", err)
	}

	TestRepo = NewTaskRepository(TestDB)

	if err = runMigrations(); err != nil {
		fatal(t, "failed to run migrations: %v", err)
	}

}

func runMigrations() error {
	err := TestDB.Exec(`
        DO $$ 
        BEGIN 
            IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'task_status') THEN
                CREATE TYPE task_status AS ENUM ('todo', 'in_progress', 'done');
            END IF;
        END $$;
    `).Error
	if err != nil {
		return fmt.Errorf("failed to create enum task_status: %w", err)
	}

	err = TestDB.AutoMigrate(
		&model.Task{},
		&model.User{},
		&model.Board{},
		&model.BoardColumn{},
		&model.Sprint{},
		&model.Group{},
	)
	if err != nil {
		return fmt.Errorf("AutoMigrate failed: %w", err)
	}
	return nil
}

func TeardownSuite(t *testing.T) {
	if t != nil {
		t.Helper()
	}

	if pgContainer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := pgContainer.Terminate(ctx)
		if err != nil {
			fmt.Printf("Warning: failed to terminate postgres container: %v\n", err)
		}
	}
}

func fatal(t *testing.T, format string, args ...interface{}) {
	if t != nil {
		t.Fatalf(format, args...)
	} else {
		panic(fmt.Sprintf(format, args...))
	}
}
