package service

import (
	"context"
	"golibrary/internal/model"
	"golibrary/internal/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Service struct {
	User UserWorker
}

type UserWorker interface {
	ListUsers(ctx context.Context) ([]model.User, error)
	GetUserByID(ctx context.Context, userId int) (*model.User, error)
	GetUserBySurname(ctx context.Context, surname string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) (int, error)
	UpdateUser(ctx context.Context, userId int, user *model.UserUpdate) error
	CreateFriendship(ctx context.Context, firstID, secondID int) error
	ListFriendships(ctx context.Context, userId int) ([]model.User, error)
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		User: NewUsersService(repos.User),
	}
}
