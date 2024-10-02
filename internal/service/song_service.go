package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"music-test-lib/internal/domain"
	"music-test-lib/internal/repository"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// SongService содержит бизнес-логику для работы с песнями.
type SongService struct {
	repo *repository.SongRepository
	log  *slog.Logger
}

// NewSongService создаёт новый экземпляр SongService.
func NewSongService(repo *repository.SongRepository, log *slog.Logger) *SongService {
	return &SongService{repo: repo, log: log}
}

// GetSongs возвращает список песен с фильтрацией.
func (s *SongService) GetSongs(params map[string][]string) ([]domain.Song, error) {
	return s.repo.GetSongs(params)
}

// GetSongLyrics возвращает текст песни.
func (s *SongService) GetSongLyrics(id string, verse string) (string, error) {
	song, err := s.repo.GetSongByID(id)
	if err != nil {
		return "", err
	}

	// Если необходимый куплет не указан - возвращаем весь текст песни
	if verse == "" {
		return song.Lyrics, nil
	}
	// Разделяем текст песни на куплеты по двум символам новой строки
	verses := strings.Split(song.Lyrics, "\\\\n\\\\n")

	// Преобразуем параметр verse в int
	verseIndex, err := strconv.Atoi(verse)
	if err != nil || verseIndex < 1 || verseIndex > len(verses) {
		return "", errors.New("некорректный номер куплета")
	}

	// Возвращаем соответствующий куплет (нумерация начинается с 1)
	return verses[verseIndex-1], nil
}

func (s *SongService) AddSong(group, songTitle, url string) (*domain.SongWithoutID, error) {
	// Логика запроса к внешнему API для получения данных о песне
	apiUrl := fmt.Sprintf("%s?group=%s&song=%s", url, group, songTitle)
	s.log.Debug("add Song ext url:", apiUrl)
	resp, err := http.Get(apiUrl)
	if err != nil {
		s.log.Error("Ошибка запроса к внешнему API: ", err)
		return nil, fmt.Errorf("ошибка запроса к внешнему API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("внешний API вернул статус: %v", resp.StatusCode)
	}
	// Чтение данных с API
	var externalSong struct {
		ReleaseDate string `json:"releaseDate"`
		Text        string `json:"text"`
		Link        string `json:"link"`
	}

	/* Тест данных из api по описанному сваггеру */
	//externalSong.ReleaseDate = "16.07.2006"
	//externalSong.Text = "Ooh baby, don't you know I suffer?\\nOoh baby, can you hear me moan?\\nYou caught me under false pretenses\\nHow long before you let me go?\\n\\nOoh\\nYou set my soul alight\\nOoh\\nYou set my soul alight"
	//externalSong.Link = "https://www.youtube.com/watch?v=Xsp3_a-PMTw"
	//

	err = json.NewDecoder(resp.Body).Decode(&externalSong)
	s.log.Info(externalSong.ReleaseDate, externalSong.Text, externalSong.Link)
	if err != nil {
		return nil, fmt.Errorf("не удалось декодировать ответ API: %v", err)
	}
	// Парсинг даты для БД
	parseDate, err := time.Parse("02.01.2006", externalSong.ReleaseDate)
	if err != nil {
		s.log.Error("Не правильная дата: ", err)
		return nil, err
	}

	// Формируем объект песни для сохранения
	newSong := &domain.SongWithoutID{
		Group:       group,
		Title:       songTitle,
		Lyrics:      externalSong.Text,
		ReleaseDate: parseDate.Format("2006-01-02"),
		Link:        externalSong.Link,
	}

	// Сохраняем песню в базу данных
	if err := s.repo.AddSong(*newSong); err != nil {
		return nil, fmt.Errorf("ошибка сохранения песни в базе данных: %v", err)
	}

	return newSong, nil
}

// UpdateSong обновляет данные песни.
func (s *SongService) UpdateSong(id string, updates map[string]string) error {
	// Получить текущие данные песни
	song, err := s.repo.GetSongByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return repository.ErrNotFound
		}
		return err
	}

	// Обновить необходимые поля
	if group, ok := updates["group_name"]; ok {
		song.Group = group
	}
	if title, ok := updates["song_name"]; ok {
		song.Title = title
	}
	if lyrics, ok := updates["lyrics"]; ok {
		song.Lyrics = lyrics
	}
	if releaseDate, ok := updates["release_date"]; ok {
		song.ReleaseDate = releaseDate
	}
	if link, ok := updates["link"]; ok {
		song.Link = link
	}

	// Сохранить обновлённую песню в базе данных
	return s.repo.UpdateSong(*song)
}

// DeleteSong удаляет песню.
func (s *SongService) DeleteSong(id string) error {
	return s.repo.DeleteSong(id)
}
