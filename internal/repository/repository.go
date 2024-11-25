package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"golibrary/internal/model"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type Repository struct {
	User UserWorker
	Init InitWorker
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

type InitWorker interface {
	InitDB() error
	CheckEmpty() bool
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User: NewUsersRepository(db),
		Init: NewInitRepository(db),
	}
}
