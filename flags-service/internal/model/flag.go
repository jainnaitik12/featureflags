package model

import "time"

// Flag is a feature flag stored in Redis.
type Flag struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

// ToggleEvent is published to Redis pub/sub on toggle.
type ToggleEvent struct {
	FlagName  string    `json:"flag_name"`
	Action    string    `json:"action"`
	OldValue  bool      `json:"old_value"`
	NewValue  bool      `json:"new_value"`
	ChangedAt time.Time `json:"changed_at"`
}
