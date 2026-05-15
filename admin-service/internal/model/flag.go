package model

// Flag mirrors flags-service JSON.
type Flag struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}
