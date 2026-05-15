package model

import "time"

// AuditEvent is stored in PostgreSQL and returned from the API.
type AuditEvent struct {
	ID        string    `json:"id,omitempty"`
	FlagName  string    `json:"flag_name"`
	Action    string    `json:"action"`
	OldValue  bool      `json:"old_value"`
	NewValue  bool      `json:"new_value"`
	ChangedAt time.Time `json:"changed_at,omitempty"`
}
