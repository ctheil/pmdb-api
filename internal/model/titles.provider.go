package model

type ProviderRegion struct {
	Link     string     `json:"link"`
	Buy      []Provider `json:"buy"`
	Rent     []Provider `json:"rent"`
	Flatrate []Provider `json:"flatrate"`
}
type Provider struct {
	LogoPath        string `json:"logo_path"`
	ProviderID      int    `json:"provider_id"`
	ProviderName    string `json:"provider_name"`
	DisplayPriority int    `json:"display_priority"`
}
type WatchProviders struct {
	Results map[string]ProviderRegion `json:"results"`
}
