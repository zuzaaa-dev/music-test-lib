package db

import (
	"fmt"
	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log/slog"
)

// DBConfig содержит настройки подключения к базе данных.
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func NewDBConfig(host string, port string, user string, password string, name string) *Config {
	return &Config{Host: host, Port: port, User: user, Password: password, Name: name}
}

// Connect создает подключение к базе данных PostgreSQL.
func Connect(cfg *Config, log *slog.Logger) (*sqlx.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)
	log.Debug("psqlInfo: ", psqlInfo)
	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	log.Info("db connected")

	return db, nil
}

// Migrate выполняет миграции базы данных.
func Migrate(cfg *Config, fileUrl string, log *slog.Logger) error {
	m, err := migrate.New(fileUrl, fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Name, cfg.Password, cfg.Host, cfg.Port, cfg.Name))
	log.Debug("migrate: ", m)
	if err != nil {
		return err
	}
	log.Debug("migrate no error: ", m)
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
