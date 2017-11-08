package config

type Youniverse struct {
	MaxSize      int64    `json:"max_size"`
	PeersURL     string   `json:"peers_url"`
	ResourceURLs []string `json:"resource_urls"`
}
