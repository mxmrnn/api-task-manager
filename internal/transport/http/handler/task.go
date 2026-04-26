package handler

import (
	"async-api-task-manager/internal/model"
	"async-api-task-manager/internal/service"
	"async-api-task-manager/internal/transport/http/dto"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type TaskHandler struct {
	taskService service.TaskService
}

func NewTaskHandler(taskService service.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

func (h *TaskHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid task id",
		})
		return
	}

	task, err := h.taskService.GetTaskByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "task not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get task",
		})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) ListTasks(c *gin.Context) {
	var filter dto.Filter

	authorID, err := parseUUIDQuery(c, "author_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid author_id"})
		return
	}
	filter.Author = authorID

	assigneeID, err := parseUUIDQuery(c, "assignee_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assignee_id"})
		return
	}
	filter.Assignee = assigneeID

	status, err := parseStatusQuery(c, "status")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}
	filter.Status = status

	tasks, err := h.taskService.ListTasks(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	var req dto.TaskUpdateRequest

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	err = h.taskService.UpdateTask(c.Request.Context(), req, id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTaskNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		case errors.Is(err, service.ErrTitleTooShort),
			errors.Is(err, service.ErrInvalidTaskStatus):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
	}
	c.Status(http.StatusNoContent)

}

func parseUUIDQuery(c *gin.Context, key string) (*uuid.UUID, error) {
	v := c.Query(key)
	if v == "" {
		return nil, nil
	}

	id, err := uuid.Parse(v)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func parseStatusQuery(c *gin.Context, key string) (*model.TaskStatus, error) {
	v := c.Query(key)
	if v == "" {
		return nil, nil
	}

	status := model.TaskStatus(v)

	return &status, nil
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req dto.TaskCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	task, err := h.taskService.CreateTask(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTitleTooShort), errors.Is(err, service.ErrInvalidTaskStatus):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	err = h.taskService.DeleteTask(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTaskNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
	}
	c.Status(http.StatusNoContent)
}
