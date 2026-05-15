package model

import "time"

// AuditEvent is sent to audit-service.
type AuditEvent struct {
	FlagName  string    `json:"flag_name"`
	Action    string    `json:"action"`
	OldValue  bool      `json:"old_value"`
	NewValue  bool      `json:"new_value"`
	ChangedAt time.Time `json:"changed_at"`
}
