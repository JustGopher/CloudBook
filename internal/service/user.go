package service

import (
	"CloudBook/internal/domain"
	"CloudBook/internal/repository"
	"context"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	// 要考虑加密放在哪里的问题
	// 然后就是存起来
	return svc.repo.Create(ctx, u)
}
