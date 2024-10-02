package v1

import (
	"errors"
	"github.com/labstack/echo/v4"
	"log/slog"
	"music-test-lib/config"
	"music-test-lib/internal/domain"
	_ "music-test-lib/internal/domain"
	"music-test-lib/internal/repository"
	"music-test-lib/internal/service"
	"net/http"
	"strconv"
)

// Handlers содержит методы-обработчики для работы с песнями.
type Handlers struct {
	logger  *slog.Logger
	service *service.SongService
	cfg     *config.Config
}

// NewHandlers создаёт новый экземпляр Handlers с переданным логгером.
func NewHandlers(logger *slog.Logger, service *service.SongService, cfg *config.Config) *Handlers {
	return &Handlers{logger: logger, service: service, cfg: cfg}
}

type ErrorResponse struct {
	Message string `json:"message" example:"Описание ошибки"`
}
type SuccessResponse struct {
	Message string `json:"message" example:"Сообщение"`
}

// GetSongs возвращает список песен с фильтрацией по всем полям и пагинацией.
// @Summary Получить список песен
// @Description Возвращает список песен с возможностью фильтрации по всем полям (группа, название, дата выпуска, текст) и пагинацией.
// @Tags songs
// @Accept  json
// @Produce  json
// @Param group_name query string false "Название группы"
// @Param song_name query string false "Название песни"
// @Param release_date query string false "Дата релиза (формат: YYYY-MM-DD)"
// @Param lyrics query string false "Часть текста песни"
// @Param page query int false "Номер страницы"
// @Param limit query int false "Количество записей на странице"
// @Success 200 {array} domain.Song "Список песен"
// @Failure 400 {object} ErrorResponse "Некорректные параметры запроса"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /songs [get]
func (h *Handlers) GetSongs(c echo.Context) error {
	h.logger.Info("GetSongs called")

	// Получаем параметры фильтрации
	groupName := c.QueryParam("group_name")
	songName := c.QueryParam("song_name")
	releaseDate := c.QueryParam("release_date")
	lyrics := c.QueryParam("lyrics")
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	// Параметры для фильтрации
	params := map[string][]string{
		"group_name":   {groupName},
		"song_name":    {songName},
		"release_date": {releaseDate},
		"lyrics":       {lyrics},
	}
	h.logger.Info("GetSongs, params:", params)

	// Получаем список песен с фильтрацией и пагинацией из сервиса
	songs, err := h.service.GetSongs(params)
	if err != nil {
		h.logger.Error("Failed to get songs", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{"Внутренняя ошибка сервера"})
	}
	h.logger.Info("GetSongs", slog.Any("songs", songs))

	// Пагинация
	start := (page - 1) * limit
	end := start + limit
	if start > len(songs) {
		return c.JSON(http.StatusOK, []domain.Song{})
	}
	if end > len(songs) {
		end = len(songs)
	}
	if songs[start:end] == nil {
		return c.JSON(http.StatusOK, []domain.Song{})
	}

	return c.JSON(http.StatusOK, songs[start:end])
}

// GetSongText возвращает текст песни с пагинацией по куплетам.
// @Summary Получить текст песни
// @Description Возвращает текст песни с возможностью пагинации по куплетам.
// @Tags songs
// @Accept  json
// @Produce  json
// @Param id path string true "ID песни"
// @Param verse query int false "Номер куплета"
// @Success 200 {string} string "Текст песни"
// @Failure 400 {object} ErrorResponse "Некорректный ID песни"
// @Failure 404 {object} ErrorResponse "Песня не найдена"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /songs/{id} [get]
func (h *Handlers) GetSongText(c echo.Context) error {
	h.logger.Info("GetSongText called", slog.String("song_id", c.Param("id")))
	songId := c.Param("id")
	verseStr := c.QueryParam("verse")

	verse, err := strconv.Atoi(verseStr)
	if err != nil || verse < 1 {
		verse = 1
	}
	lyrics, err := h.service.GetSongLyrics(songId, verseStr)
	if err != nil {
		h.logger.Info("GetSongText Lyrics failed", slog.String("err", err.Error()))
		if errors.Is(err, repository.ErrNotFound) {
			return c.JSON(http.StatusNotFound, "Песня не найдена")
		}
		return c.JSON(http.StatusInternalServerError, "Ошибка при получении песни")
	}

	return c.JSON(http.StatusOK, lyrics)
}

type AddSongRequest struct {
	Group string `json:"group" example:"Muse"`
	Title string `json:"song" example:"Supermassive Black Hole"`
}

type ExternalAPISongResponse struct {
	ReleaseDate string `json:"releaseDate" example:"16.07.2006"`
	Text        string `json:"text" example:"Ooh baby, don't you know I suffer?\\nOoh baby, can you hear me moan?\\nYou caught me under false pretenses\\nHow long before you let me go?\\n\\nOoh\\nYou set my soul alight\\nOoh\\nYou set my soul alight"`
	Link        string `json:"link" example:"https://www.youtube.com/watch?v=Xsp3_a-PMTw"`
}

// AddSong добавляет новую песню в библиотеку.
// @Summary Добавить новую песню
// @Description Добавляет новую песню в библиотеку и получает информацию о песне из внешнего API.
// @Tags songs
// @Accept  json
// @Produce  json
// @Param song body AddSongRequest true "Песня (группа и название)"
// @Success 201 {object} domain.Song "Песня добавлена с детальной информацией"
// @Failure 400 {object} ErrorResponse "Некорректные данные песни"
// @Failure 500 {object} ErrorResponse "Не удалось добавить песню"
// @Router /songs [post]
func (h *Handlers) AddSong(c echo.Context) error {
	h.logger.Info("AddSong called")

	// Парсим запрос от клиента
	var addSongRequest AddSongRequest
	if err := c.Bind(&addSongRequest); err != nil {
		h.logger.Error("Invalid song data", err)
		return c.JSON(http.StatusBadRequest, ErrorResponse{"Некорректные данные песни"})
	}
	h.logger.Info(addSongRequest.Group, addSongRequest.Title)

	// Проверяем наличие группы и названия песни
	if addSongRequest.Group == "" || addSongRequest.Title == "" {
		h.logger.Warn("Group or song title is missing")
		return c.JSON(http.StatusBadRequest, ErrorResponse{"Группа и название песни обязательны"})
	}

	// Вызываем метод сервиса для добавления песни
	newSong, err := h.service.AddSong(addSongRequest.Group, addSongRequest.Title, h.cfg.API.MusicInfoURL)
	if err != nil {
		h.logger.Error("Failed to add song", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{"Не удалось добавить песню"})
	}

	// Возвращаем добавленную песню
	return c.JSON(http.StatusCreated, newSong)
}

type UpdateSongRequest struct {
	Group       string `json:"group" example:"Muse"`
	Title       string `json:"title" example:"Supermassive Black Hole"`
	ReleaseDate string `json:"release_date" example:"2006-07-16"`
	Lyrics      string `json:"lyrics" example:"Ooh baby, don't you know I suffer? ..."`
	Link        string `json:"link" example:"https://www.youtube.com/watch?v=Xsp3_a-PMTw"`
}

// UpdateSong обновляет данные существующей песни.
// @Summary Изменить данные песни
// @Description Обновляет информацию о песне, включая название группы, название песни, текст, дату релиза и ссылку.
// @Tags songs
// @Accept  json
// @Produce  json
// @Param id path string true "ID песни"
// @Param song body UpdateSongRequest true "Новая информация о песне"
// @Success 200 {object} SuccessResponse "Данные песни обновлены"
// @Failure 400 {object} ErrorResponse "Некорректные данные песни"
// @Failure 404 {object} ErrorResponse "Песня не найдена"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /songs/{id} [put]
func (h *Handlers) UpdateSong(c echo.Context) error {
	h.logger.Info("UpdateSong called", slog.String("song_id", c.Param("id")))

	// Получаем ID песни из параметров пути
	id := c.Param("id")

	// Парсинг данных обновления из тела запроса
	var updateReq UpdateSongRequest
	if err := c.Bind(&updateReq); err != nil {
		h.logger.Error("Invalid song data", slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, ErrorResponse{"Некорректные данные песни"})
	}

	// Подготовка данных для обновления
	updates := map[string]string{}
	if updateReq.Group != "" {
		updates["group_name"] = updateReq.Group
	}
	if updateReq.Title != "" {
		updates["song_name"] = updateReq.Title
	}
	if updateReq.Lyrics != "" {
		updates["lyrics"] = updateReq.Lyrics
	}
	if updateReq.ReleaseDate != "" {
		updates["release_date"] = updateReq.ReleaseDate
	}
	if updateReq.Link != "" {
		updates["link"] = updateReq.Link
	}

	// Обновляем песню через сервис
	err := h.service.UpdateSong(id, updates)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.logger.Warn("Song not found", slog.String("song_id", id))
			return c.JSON(http.StatusNotFound, ErrorResponse{"Песня не найдена"})
		}
		h.logger.Error("Failed to update song", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{"Не удалось обновить песню"})
	}

	// Возвращаем успешный ответ
	return c.JSON(http.StatusOK, SuccessResponse{"Данные песни обновлены"})
}

// DeleteSong удаляет песню из библиотеки.
// @Summary Удалить песню
// @Description Удаляет песню по ее ID.
// @Tags songs
// @Param id path string true "ID песни"
// @Success 200 {object} SuccessResponse "Песня удалена"
// @Failure 404 {object} ErrorResponse "Песня не найдена"
// @Failure 500 {object} ErrorResponse "Не удалось удалить песню"
// @Router /songs/{id} [delete]
func (h *Handlers) DeleteSong(c echo.Context) error {
	h.logger.Info("DeleteSong called", slog.String("song_id", c.Param("id")))

	// Получаем ID песни из параметров пути
	id := c.Param("id")

	// Попытка удаления песни через сервис
	err := h.service.DeleteSong(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.logger.Warn("Song not found", slog.String("song_id", id))
			return c.JSON(http.StatusNotFound, ErrorResponse{"Песня не найдена"})
		}
		h.logger.Error("Ошибка при удалении песни", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{"Не удалось удалить песню"})
	}

	// Возвращаем успешный ответ
	h.logger.Info("Song deleted successfully", slog.String("song_id", id))
	return c.JSON(http.StatusOK, SuccessResponse{"Песня успешно удалена"})
}
