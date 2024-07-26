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

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ProductionCompany struct {
	ID            int    `json:"id"`
	LogoPath      string `json:"logo_path"`
	OriginCountry string `json:"origin_country"`
}

type Video struct {
	ISOLang639  string `json:"iso_639_1"`
	ISOLang3166 string `json:"iso_3166_1"`
	Name        string `json:"name"`
	Key         string `json:"key"`
	Site        string `json:"site"`
	Size        int    `json:"size"`
	Type        string `json:"type"`
	Official    bool   `json:"official"`
	PublishedAt string `json:"published_at"`
	ID          string `json:"id"`
}

type Credit struct {
	ProfilePath string `json:"profile_path"`
	// Gender      string `json:"gender"`
	KnownForDepartment string `json:"known_for_department"`
	Name               string `json:"name"`
	// OriginalName       string  `json:"original_name"`
	// Popularity         float32 `json:"popularity"`
	Character string `json:"character"`
	// Order              int     `json:"order"`
	// ID                 int     `json:"id"`
	// CastID             int     `json:"cast_id"`
	// CreditID           int     `json:"credit_id"`
	// Adult              bool    `json:"adult"`
}
type Credits struct {
	ID   int      `json:"id"`
	Cast []Credit `json:"cast"`
	// Crew []Credit `json:"crew"`
}
type Videos struct {
	ID      int     `json:"id"`
	Results []Video `json:"results"`
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
}
