package service

import (
	"context"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/application/dto"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/domain"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/domain/entity"
)

type UserService interface {
	GetAllUsers(ctx context.Context) ([]dto.UserOutput, error)
	CreateUser(ctx context.Context, user entity.User) (entity.User, error)
}

type userService struct {
	userRepo domain.UserRepository
}

func NewUserService(userRepo domain.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetAllUsers(ctx context.Context) ([]dto.UserOutput, error) {
	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	userOutputs := make([]dto.UserOutput, 0, len(users))
	for _, user := range users {
		if user.ID == 0 {
			continue
		}

		userOutput := dto.UserOutput{
			Name:  user.Name,
			Email: user.Email,
			Age:   user.Age,
		}

		userOutputs = append(userOutputs, userOutput)
	}

	return userOutputs, nil
}

func (s *userService) CreateUser(ctx context.Context, user entity.User) (entity.User, error) {
	return s.userRepo.Create(ctx, user)
}
