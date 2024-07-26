package model

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
type Videos struct {
	ID      int     `json:"id"`
	Results []Video `json:"results"`
}
