package config

type Redirect struct {
	PAC          bool   `json:"pac"`
	Encode       bool   `json:"encode"`
	RulesURL     string `json:"rules_url"`
	BricksURL    string `json:"bricks_url"`
	UpstreamsURL string `json:"services_url"`
}
