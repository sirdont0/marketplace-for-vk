package storage

import (
	"context"
	"errors"

	"log"
	"marketplace/internal/models"

	"github.com/jackc/pgconn"
)

func CreateUser(user *models.User) error {
	query := `INSERT INTO users (login, password_hash) VALUES ($1, $2) RETURNING id, created_at`
	err := db.QueryRow(context.Background(), query, user.Login, user.PasswordHash).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		// Проверяем на ошибку уникальности
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return errors.New("user already exists")
		}
		// Логируем и возвращаем другую ошибку
		log.Printf("CreateUser error: %v", err)
		return errors.New("internal server error")
	}
	return nil
}

func GetUserByLogin(login string) (*models.User, error) {
	var user models.User
	query := `SELECT id, login, password_hash, created_at FROM users WHERE login = $1`
	err := db.QueryRow(context.Background(), query, login).Scan(&user.ID, &user.Login, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}
