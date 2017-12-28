package config

type Redirect struct {
	Encode       bool   `json:"encode"`
	RulesURL     string `json:"rules_url"`
	UpstreamsURL string `json:"services_url"`
}
