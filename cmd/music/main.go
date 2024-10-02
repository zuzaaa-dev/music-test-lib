package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"
	"log"
	"log/slog"
	"music-test-lib/config"
	_ "music-test-lib/docs"
	v1 "music-test-lib/internal/api/v1"
	"music-test-lib/internal/repository"
	"music-test-lib/internal/service"
	"music-test-lib/pkg/db"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// @title Online Music Library API
// @version 1.0
// @description This is an API for an online music library, providing functionality to manage and query songs.
// @host localhost:8080
func main() {

	// Загружаем переменные из .env файла
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	cfg := config.MustLoad()

	// Настраиваем логгер
	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))

	// Логируем инициализацию сервера
	log.Info("initializing server", slog.String("address", cfg.HTTPServer.Address))
	log.Debug("logger debug mode enabled")

	dbConfig := &db.Config{
		Host:     cfg.DataBase.Host,
		Port:     cfg.DataBase.Port,
		User:     cfg.DataBase.User,
		Password: cfg.DataBase.Password,
		Name:     cfg.DataBase.Name,
	}

	dbConn, err := db.Connect(dbConfig, log)
	if err != nil {
		log.Error("failed to connect to database: %v", err)
	}
	defer dbConn.Close()
	log.Info("connect db success")

	makeMigrate(dbConfig, cfg.DataBase.FileMigrations, log)

	repo := repository.NewSongRepository(dbConn)
	songService := service.NewSongService(repo, log)

	e := echo.New()

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	v1.RegisterRoutes(e, log, songService, cfg)

	e.Logger.Fatal(e.Start(cfg.HTTPServer.Address))
}

func makeMigrate(cfg *db.Config, filePath string, log *slog.Logger) {
	if err := db.Migrate(cfg, filePath, log); err != nil {
		log.Error("failed to run migrations: ", err)
		os.Exit(1)
	} else {
		log.Info("migration success")
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
