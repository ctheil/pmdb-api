package model

type Title struct {
	OriginalTitle    string   `json:"original_title"`
	Overview         string   `json:"overview"`
	Genres           []string `json:"genres"`
	OriginalLanguage string   `json:"original_language"`
	PosterPath       string   `json:"poster_path"`
	Title            string   `json:"title"`
	Status           string   `json:"status"`
	IMDB_ID          string   `json:"imdb_id"`
	Id               int      `json:"id"`
	ReleaseDate      string   `json:"release_date"`
	Runtime          int      `json:"runtime"`
}
