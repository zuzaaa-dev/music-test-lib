package domain

// Song представляет структуру песни.
// @Description Модель данных песни.
//
//	@Example {
//	  "id": "1",
//	  "group": "Muse",
//	  "title": "Supermassive Black Hole",
//	  "release_date": "2006-07-16",
//	  "lyrics": "Ooh baby, don't you know I suffer? ...",
//	  "link": "https://www.youtube.com/watch?v=Xsp3_a-PMTw"
//	}
type Song struct {
	ID          string `json:"id" example:"1"`
	Group       string `json:"group" example:"Muse"`
	Title       string `json:"title" example:"Supermassive Black Hole"`
	ReleaseDate string `json:"release_date" example:"2006-07-16"`
	Lyrics      string `json:"lyrics" example:"Ooh baby, don't you know I suffer? ..."`
	Link        string `json:"link" example:"https://www.youtube.com/watch?v=Xsp3_a-PMTw"`
}

// SongWithoutID представляет структуру песни без ID.
// @Description Модель данных песни.
//
//	@Example {
//	  "group": "Muse",
//	  "title": "Supermassive Black Hole",
//	  "release_date": "2006-07-16",
//	  "lyrics": "Ooh baby, don't you know I suffer? ...",
//	  "link": "https://www.youtube.com/watch?v=Xsp3_a-PMTw"
//	}
type SongWithoutID struct {
	Group       string `json:"group" example:"Muse"`
	Title       string `json:"title" example:"Supermassive Black Hole"`
	ReleaseDate string `json:"release_date" example:"2006-07-16"`
	Lyrics      string `json:"lyrics" example:"Ooh baby, don't you know I suffer? ..."`
	Link        string `json:"link" example:"https://www.youtube.com/watch?v=Xsp3_a-PMTw"`
}
