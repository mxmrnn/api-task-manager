package rpc

import (
	"async-api-task-manager/internal/service"
	"async-api-task-manager/internal/transport/http/dto"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

type Handler struct {
	service service.TaskService
}

func NewHandler(svc service.TaskService) *Handler {
	return &Handler{service: svc}
}

func (h *Handler) GetUserTaskStatsHandler(payload []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := dto.TaskStatsRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, fmt.Errorf("invalid payload: %w", err)
	}

	if req.UserID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	id, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id uuid format: %w", err)
	}

	stats, err := h.service.GetUserTaskCounts(ctx, id)
	if err != nil {
		return json.Marshal(map[string]string{
			"code":    "internal error",
			"message": err.Error(),
		})
	}

	response := dto.TaskStatsResponse{
		AsAssignee: stats.AsAssignee,
		AsWatcher:  stats.AsWatcher,
	}

	return json.Marshal(response)
}
