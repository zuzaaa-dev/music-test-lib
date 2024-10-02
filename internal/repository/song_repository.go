package repository

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"music-test-lib/internal/domain"
)

var ErrNotFound = errors.New("song not found")

// SongRepository определяет интерфейс для работы с песнями в базе данных.
type SongRepository struct {
	db *sqlx.DB
}

// NewSongRepository создает новый SongRepository.
func NewSongRepository(db *sqlx.DB) *SongRepository {
	return &SongRepository{db: db}
}

// GetSongs возвращает список песен с фильтрацией и пагинацией.
func (r *SongRepository) GetSongs(params map[string][]string) ([]domain.Song, error) {
	var songs []domain.Song
	query := "SELECT id, group_name, song_name, lyrics, release_date, link FROM songs WHERE 1=1"

	// Фильтрация по полям
	if groupName, ok := params["group_name"]; ok && groupName[0] != "" {
		query += " AND group_name ILIKE '%" + groupName[0] + "%'"
	}
	if songName, ok := params["song_name"]; ok && songName[0] != "" {
		query += " AND song_name ILIKE '%" + songName[0] + "%'"
	}
	if releaseDate, ok := params["release_date"]; ok && releaseDate[0] != "" {
		query += " AND release_date = '" + releaseDate[0] + "'"
	}
	if lyrics, ok := params["lyrics"]; ok && lyrics[0] != "" {
		query += " AND lyrics ILIKE '%" + lyrics[0] + "%'"
	}

	// Выполнение запроса
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Чтение результатов
	for rows.Next() {
		var song domain.Song
		if err := rows.Scan(
			&song.ID,
			&song.Group,
			&song.Title,
			&song.Lyrics,
			&song.ReleaseDate,
			&song.Link,
		); err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}
	return songs, nil
}

// GetSongByID возвращает текст песни по её ID.
func (r *SongRepository) GetSongByID(id string) (*domain.Song, error) {
	var song domain.Song
	query := "SELECT id, group_name, song_name, release_date, lyrics, link FROM songs WHERE id = $1"
	err := r.db.QueryRow(query, id).Scan(&song.ID, &song.Group, &song.Title, &song.ReleaseDate, &song.Lyrics, &song.Link)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return &song, nil
}

// AddSong добавляет новую песню в базу данных.
func (r *SongRepository) AddSong(song domain.SongWithoutID) error {
	_, err := r.db.Exec(
		"INSERT INTO songs (group_name, song_name, release_date, lyrics, link) VALUES ($1, $2, $3, $4, $5)",
		song.Group, song.Title, song.ReleaseDate, song.Lyrics, song.Link,
	)

	return err
}

// UpdateSong обновляет данные песни.
func (r *SongRepository) UpdateSong(song domain.Song) error {
	_, err := r.db.Exec(
		"UPDATE songs SET group_name= $2, song_name = $3, release_date = $4, lyrics = $5, link = $6 WHERE id = $1",
		song.ID, song.Group, song.Title, song.ReleaseDate, song.Lyrics, song.Link,
	)
	return err
}

// DeleteSong удаляет песню из базы данных.
func (r *SongRepository) DeleteSong(id string) error {
	_, err := r.db.Exec("DELETE FROM songs WHERE id = $1", id)
	return err
}
