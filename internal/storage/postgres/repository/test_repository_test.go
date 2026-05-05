//task_integration_test.go
//go:build integration

package repository_test

import (
	"async-api-task-manager/internal/storage/postgres/repository"
	taskdomain "async-api-task-manager/internal/transport/http/dto"
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"async-api-task-manager/internal/model"
)

func TestMain(m *testing.M) {
	repository.SetupSuite(nil)

	code := m.Run()

	repository.TeardownSuite(nil)

	os.Exit(code)
}

func TestTaskRepository_CreateAndGet(t *testing.T) {

	repo := repository.TestRepo
	author := createTestAuthor(t)

	task := &model.Task{
		Title:       "test create",
		Description: "test description",
		Status:      model.TaskStatusTodo,
		AuthorID:    author.ID,
		AssigneeID:  nil,
		BoardID:     nil,
		BoardColumn: nil,
		Group:       nil,
		Sprint:      nil,
	}

	err := repo.Create(context.Background(), task)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, task.ID)

	found, err := repo.GetByID(context.Background(), task.ID)
	require.NoError(t, err)
	assert.Equal(t, task.Title, found.Title)
	assert.NotNil(t, found.Author)
}

func TestTaskRepository_List(t *testing.T) {
	repo := repository.TestRepo
	author := createTestAuthor(t)

	createTestTask(t, repo, author.ID, "task 1", model.TaskStatusTodo)
	createTestTask(t, repo, author.ID, "task 2", model.TaskStatusDone)

	filter := taskdomain.Filter{Status: ptr(model.TaskStatusTodo)}

	tasks, err := repo.List(context.Background(), filter)

	require.NoError(t, err)
	require.NotEmpty(t, tasks)
}

func TestTaskRepository_Update(t *testing.T) {
	repo := repository.TestRepo
	author := createTestAuthor(t)

	task := &model.Task{
		Title:    "task for update",
		Status:   model.TaskStatusTodo,
		AuthorID: author.ID,
	}
	err := repo.Create(context.Background(), task)
	require.NoError(t, err)

	task.Title = "Обновлённый заголовок"
	task.Status = model.TaskStatusInProgress

	err = repo.Update(context.Background(), task)
	require.NoError(t, err)

	updated, err := repo.GetByID(context.Background(), task.ID)
	require.NoError(t, err)
	assert.Equal(t, "Обновлённый заголовок", updated.Title)
	assert.Equal(t, model.TaskStatusInProgress, updated.Status)
}

func TestTaskRepository_Delete(t *testing.T) {
	repo := repository.TestRepo
	author := createTestAuthor(t)

	task := &model.Task{
		Title:    "Task for delete",
		AuthorID: author.ID,
	}
	err := repo.Create(context.Background(), task)
	require.NoError(t, err)

	err = repo.Delete(context.Background(), task.ID)
	require.NoError(t, err)

	_, err = repo.GetByID(context.Background(), task.ID)
	assert.ErrorIs(t, err, repository.ErrTaskNotFound)
}

func createTestAuthor(t *testing.T) *model.User {
	t.Helper()

	author := &model.User{
		FullName: "test author name",
		Email:    "author." + uuid.NewString()[:8] + "@example.com",
	}
	require.NoError(t, repository.TestDB.Create(author).Error)
	return author
}

func createTestTask(t *testing.T, repo *repository.TaskRepository, authorID uuid.UUID, title string, status model.TaskStatus) *model.Task {
	t.Helper()
	task := &model.Task{
		Title:    title,
		Status:   status,
		AuthorID: authorID,
	}
	require.NoError(t, repo.Create(context.Background(), task))
	return task
}

func ptr[T any](v T) *T { return &v }
