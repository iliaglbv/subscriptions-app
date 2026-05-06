// internal/repository/user_repo.go
package repository

import (
	"database/sql"
	"errors"
	"time"

	"subscriptions-app/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(username, email, passwordHash string) (int64, error) {
	var id int64
	err := r.db.QueryRow(
		`INSERT INTO users (username, email, password_hash, created_at) 
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		username, email, passwordHash, time.Now(),
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(
		`SELECT id, username, email, password_hash, created_at 
		 FROM users WHERE username = $1`,
		username,
	).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(
		`SELECT id, username, email, created_at FROM users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
