package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type InitRepository struct {
	db *sqlx.DB
}

func NewInitRepository(db *sqlx.DB) *InitRepository {
	return &InitRepository{db: db}
}

func (r *InitRepository) InitDB() error {
	query := fmt.Sprint(`
CREATE TABLE IF NOT EXISTS users(
    id serial unique not null,
    name varchar(50) not null,
    surname varchar(50) not null,
    middle_name varchar(50) not null,
    age integer not null,
    nation varchar(50) not null,
    gender varchar(10) not null
);

CREATE TABLE IF NOT EXISTS emails(
    id serial unique not null,
    email varchar(50) unique not null,
    user_id integer references users(id) on delete cascade
);

CREATE TABLE IF NOT EXISTS friendships(
  user_id1 integer references users(id) on delete cascade,
  user_id2 integer references users(id) on delete cascade,
  primary key (user_id1, user_id2)
);`)

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (r *InitRepository) CheckEmpty() bool {
	query := fmt.Sprintf("SELECT * FROM users")

	_, err := r.db.Query(query)
	if err != nil {
		return true
	}

	return false
}
