package repository

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"golibrary/internal/model"
)

type UsersRepository struct {
	db *sqlx.DB
}

func NewUsersRepository(db *sqlx.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (r *UsersRepository) ListUsers(ctx context.Context) ([]model.User, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("Error to begin tx on ListUser: %w", err)
	}

	var users []model.User
	query := fmt.Sprintf("SELECT * FROM %s", userTable)
	err = r.db.SelectContext(ctx, &users, query)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("Error to list users: %w", err)
	}

	err = r.GetEmails(ctx, users)

	return users, tx.Commit()
}

func (r *UsersRepository) GetUserByID(ctx context.Context, userId int) (*model.User, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("Error to begin tx on GetUserByID:", err)
	}

	var user model.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", userTable)
	err = r.db.GetContext(ctx, &user, query, userId)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("Error to get user by id on repo: %w", err)
	}

	var emails []model.Email

	query = fmt.Sprintf("SELECT email FROM %s WHERE user_id = $1", emailsTable)
	err = r.db.SelectContext(ctx, &emails, query, userId)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("Error to get user emails by id on repo: %w", err)
	}

	for _, email := range emails {
		user.Emails = append(user.Emails, email.Email)
	}

	return &user, nil
}

func (r *UsersRepository) GetUserBySurname(ctx context.Context, surname string) (*model.User, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("Error to begin tx on GetUserByID:", err)
	}

	var user model.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE surname = $1", userTable)
	err = r.db.GetContext(ctx, &user, query, surname)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("Error to get user by id on repo: %w", err)
	}

	var emails []model.Email

	query = fmt.Sprintf("SELECT email FROM %s WHERE id = $1", emailsTable)
	err = r.db.SelectContext(ctx, &emails, query, user.Id)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("Error to get user emails by id on repo: %w", err)
	}

	for _, email := range emails {
		user.Emails = append(user.Emails, email.Email)
	}

	return &user, nil
}

func (r *UsersRepository) CreateUser(ctx context.Context, user *model.User) (int, error) {
	var id int
	tx, err := r.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("Error to begin tx on CreateUser: %w", err)
	}

	query := fmt.Sprintf("INSERT INTO %s (name, surname, middle_name, age, nation, gender) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", userTable)
	row := r.db.QueryRowContext(ctx, query, user.Name, user.Surname, user.MiddleName, user.Age, user.Nationality, user.Gender)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("Error creating user: %w", err)
	}

	for _, email := range user.Emails {
		query = fmt.Sprintf("INSERT INTO %s (email, user_id) VALUES($1, $2)", emailsTable)
		_, err := r.db.ExecContext(ctx, query, email, id)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("Error inserting email: %w", err)
		}
	}

	return id, tx.Commit()
}

func (r *UsersRepository) UpdateUser(ctx context.Context, userId int, user *model.UserUpdate) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("Error to begin tx on UpdateUser:", err)
	}

	query := fmt.Sprintf("UPDATE %s SET name = $1, surname = $2, middle_name = $3, age = $4, nation = $5, gender = $6 WHERE id = $7", userTable)
	_, err = r.db.ExecContext(ctx, query, user.Name, user.Surname, user.MiddleName, user.Age, user.Nationality, user.Gender, userId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Error updating user: %w", err)
	}

	return tx.Commit()
}

func (r *UsersRepository) CreateFriendship(ctx context.Context, firstID, secondID int) error {
	query := fmt.Sprintf("INSERT INTO %s (user_id1, user_id2) VALUES($1, $2), ($2, $1) ON CONFLICT DO NOTHING", friendsTable)
	_, err := r.db.ExecContext(ctx, query, firstID, secondID)
	if err != nil {
		return fmt.Errorf("Error creating friendship: %w", err)
	}

	return nil
}

func (r *UsersRepository) ListFriendships(ctx context.Context, userId int) ([]model.User, error) {
	tx, err := r.db.Begin()
	var users []model.User

	query := fmt.Sprint("SELECT users.* FROM friendships JOIN users ON users.id = friendships.user_id2 WHERE friendships.user_id1 = $1")

	err = r.db.SelectContext(ctx, &users, query, userId)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("Error listing friendships: %w", err)
	}

	err = r.GetEmails(ctx, users)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("Error getting friendships: %w", err)
	}

	return users, tx.Commit()
}

func (r *UsersRepository) GetEmails(ctx context.Context, users []model.User) error {
	var emails []model.Email

	for id, user := range users {
		query := fmt.Sprintf("SELECT email FROM %s WHERE user_id = $1", emailsTable)
		err := r.db.SelectContext(ctx, &emails, query, user.Id)
		if err != nil {
			return fmt.Errorf("Error getting emails: %w", err)
		}

		for _, email := range emails {
			users[id].Emails = append(users[id].Emails, email.Email)
		}
	}

	return nil
}
