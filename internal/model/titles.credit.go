package model

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
