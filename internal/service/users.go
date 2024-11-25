package service

import (
	"context"
	"golibrary/internal/model"
	"golibrary/internal/repository"
)

type UsersService struct {
	repo repository.UserWorker
}

func NewUsersService(repo repository.UserWorker) *UsersService {
	return &UsersService{repo: repo}
}

func (p *UsersService) ListUsers(ctx context.Context) ([]model.User, error) {
	return p.repo.ListUsers(ctx)
}

func (p *UsersService) GetUserByID(ctx context.Context, userId int) (*model.User, error) {
	return p.repo.GetUserByID(ctx, userId)
}

func (p *UsersService) GetUserBySurname(ctx context.Context, surname string) (*model.User, error) {
	return p.repo.GetUserBySurname(ctx, surname)
}
func (p *UsersService) CreateUser(ctx context.Context, user *model.User) (int, error) {
	return p.repo.CreateUser(ctx, user)
}
func (p *UsersService) UpdateUser(ctx context.Context, userId int, user *model.UserUpdate) error {
	return p.repo.UpdateUser(ctx, userId, user)
}

func (p *UsersService) CreateFriendship(ctx context.Context, firstID, secondID int) error {
	return p.repo.CreateFriendship(ctx, firstID, secondID)
}

func (p *UsersService) ListFriendships(ctx context.Context, userId int) ([]model.User, error) {
	return p.repo.ListFriendships(ctx, userId)
}
