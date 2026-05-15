package model

// Flag is a feature flag returned by the public SDK API.
type Flag struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

// FlagEvaluation is the JSON shape for GET /flags/:name.
type FlagEvaluation struct {
	Flag    string `json:"flag"`
	Enabled bool   `json:"enabled"`
}
