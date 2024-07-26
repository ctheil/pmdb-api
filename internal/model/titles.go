package model

type Title struct {
	OriginalTitle    string   `json:"original_title"`
	Overview         string   `json:"overview"`
	Genres           []string `json:"genres"`
	OriginalLanguage string   `json:"original_language"`
	PosterPath       string   `json:"poster_path"`
	BackdropPath     string   `json:"backdrop_path"`
	Title            string   `json:"title"`
	Status           string   `json:"status"`
	IMDB_ID          string   `json:"imdb_id"`
	Id               int      `json:"id"`
	ReleaseDate      string   `json:"release_date"`
	Runtime          int      `json:"runtime"`
}

type TitleResponse struct {
	Page    int     `json:"page"`
	Results []Title `json:"results"`
}

type TitleDetails struct {
	BelongsToCollection struct {
		ID           int    `json:"id"`
		Name         string `json:"name"`
		PosterPath   string `json:"poster_path"`
		BackdropPath string `json:"backdrop_path"`
	} `json:"belongs_to_collection"`

	Budget              int                 `json:"budget"`
	Genres              []Genre             `json:"genres"`
	HomePage            string              `json:"homepage"`
	ID                  int                 `json:"id"`
	IMDB_ID             string              `json:"imdb_id"`
	Popularity          float64             `json:"popularity"`
	ProductionCompanies []ProductionCompany `json:"production_companies"`
	Revenue             int                 `json:"revenue"`
	Runtime             int                 `json:"runtime"`
	Status              string              `json:"status"`
	Tagline             string              `json:"tagline"`
	VoteAverage         float64             `json:"vote_average"`
	VoteCount           int                 `json:"vote_count"`
	Videos              Videos              `json:"videos"`
	Credits             Credits             `json:"credits"`
	WatchProviders      WatchProviders      `json:"watch/providers"`
}
