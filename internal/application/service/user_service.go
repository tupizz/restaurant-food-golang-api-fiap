package service

import (
	"context"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/domain"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/domain/entity"
)

type UserService interface {
	GetAllUsers(ctx context.Context) ([]entity.User, error)
	CreateUser(ctx context.Context, user entity.User) (entity.User, error)
}

type userService struct {
	userRepo domain.UserRepository
}

func NewUserService(userRepo domain.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	return s.userRepo.GetAll(ctx)
}

func (s *userService) CreateUser(ctx context.Context, user entity.User) (entity.User, error) {
	return s.userRepo.Create(ctx, user)
}
