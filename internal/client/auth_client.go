package client

import (
	"context"
	"encoding/json"
	"fmt"

	"async-api-task-manager/internal/transport/http/dto"
	"github.com/google/uuid"
)

type AuthClient interface {
	GetUserInfo(ctx context.Context, userID uuid.UUID) (*dto.UserShort, error)
}

type authClient struct {
	rpc *RPCClient
}

func NewAuthClient(rpcClient *RPCClient) AuthClient {
	return &authClient{
		rpc: rpcClient,
	}
}

func (c *authClient) GetUserInfo(ctx context.Context, userID uuid.UUID) (*dto.UserShort, error) {
	payload := map[string]string{
		"user_id": userID.String(),
	}

	respBody, err := c.rpc.Call(ctx, "auth-service.rpc", payload)
	if err != nil {
		return nil, fmt.Errorf("auth-service rpc call failed: %w", err)
	}

	var user dto.UserShort
	if err := json.Unmarshal(respBody, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user info: %w", err)
	}

	return &user, nil
}
