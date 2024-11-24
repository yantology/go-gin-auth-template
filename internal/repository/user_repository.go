package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/yantology/go-gin-auth-template/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	log.Println("Creating user repository")
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query, user.Username, user.Email, user.PasswordHash).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	fmt.Println("create query")
	query := `
		SELECT id, username, email, password_hash, created_at, updated_at
		FROM users WHERE email = $1`
	fmt.Println("do create query")
	err := r.db.QueryRow(query, email).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash,
			&user.CreatedAt, &user.UpdatedAt)

	fmt.Println("finish create query")
	if err == sql.ErrNoRows {
		return nil, nil
	}

	fmt.Println("return create query")
	return user, err
}

func (r *UserRepository) UpdatePassword(userID int, passwordHash string) error {
	query := `
		UPDATE users 
		SET password_hash = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2`

	_, err := r.db.Exec(query, passwordHash, userID)
	return err
}

func (r *UserRepository) GetById(id int) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, username, email, password_hash, created_at, updated_at
		FROM users WHERE id = $1`

	err := r.db.QueryRow(query, id).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash,
			&user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}
