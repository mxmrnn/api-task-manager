package service

import (
	"async-api-task-manager/internal/model"
	"async-api-task-manager/internal/transport/http/dto"
	"context"
	"strings"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
}

type UserService interface {
	CreateUser(ctx context.Context, req dto.UserCreateRequest) (dto.UserCreatedResponse, error)
}

type userService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (u *userService) CreateUser(ctx context.Context, req dto.UserCreateRequest) (dto.UserCreatedResponse, error) {
	us := model.User{
		FullName: strings.TrimSpace(req.FullName),
		Email:    strings.TrimSpace(req.Email),
	}

	if us.FullName == "" {
		return dto.UserCreatedResponse{}, ErrUserFullNameRequired
	}
	if us.Email == "" {
		return dto.UserCreatedResponse{}, ErrUserEmailRequired
	}

	err := u.userRepo.Create(ctx, &us)
	if err != nil {
		return dto.UserCreatedResponse{}, err
	}

	res := dto.UserCreatedResponse{
		ID:        us.ID,
		Email:     us.Email,
		FullName:  us.FullName,
		CreatedAt: us.CreatedAt,
	}

	return res, nil
}
