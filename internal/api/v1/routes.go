package v1

import (
	"github.com/labstack/echo/v4"
	"log/slog"
	"music-test-lib/config"
	"music-test-lib/internal/service"
)

// RegisterRoutes регистрирует маршруты для работы с API песен.
func RegisterRoutes(e *echo.Echo, logger *slog.Logger, service *service.SongService, cfg *config.Config) {
	handlers := NewHandlers(logger, service, cfg)

	// REST методы для библиотеки песен
	e.GET("/songs", handlers.GetSongs)          // Получение списка песен с фильтрацией и пагинацией
	e.GET("/songs/:id", handlers.GetSongText)   // Получение текста песни с пагинацией по куплетам
	e.POST("/songs", handlers.AddSong)          // Добавление новой песни
	e.PUT("/songs/:id", handlers.UpdateSong)    // Изменение данных песни
	e.DELETE("/songs/:id", handlers.DeleteSong) // Удаление песни
}
