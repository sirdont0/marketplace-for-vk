package storage

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var db *pgxpool.Pool

func InitDB() error {
	// Загружаем переменные окружения
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:5432/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return fmt.Errorf("unable to parse database config: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %v", err)
	}

	db = pool

	// Автоматическое создание таблиц, если их нет
	createTables := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			login VARCHAR(50) UNIQUE NOT NULL,
			password_hash VARCHAR(100) NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS ads (
			id SERIAL PRIMARY KEY,
			title VARCHAR(100) NOT NULL,
			text TEXT NOT NULL,
			image_url VARCHAR(255) NOT NULL,
			price DECIMAL(10, 2) NOT NULL,
			author_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMP DEFAULT NOW()
		);`,
	}

	for _, query := range createTables {
		_, err := db.Exec(context.Background(), query)
		if err != nil {
			return fmt.Errorf("failed to create table: %v", err)
		}
	}
	return nil
}
