package repository

import (
	"net/http"

	"featureflags/admin-service/internal/client"
)

// AuditUpstream calls audit-service.
type AuditUpstream struct {
	baseURL string
}

func NewAuditUpstream(baseURL string) *AuditUpstream {
	return &AuditUpstream{baseURL: baseURL}
}

func (u *AuditUpstream) Log(body []byte) (int, []byte, error) {
	return client.Do(http.MethodPost, u.baseURL+"/audit", body)
}
