package service

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"featureflags/admin-service/internal/model"
	"featureflags/admin-service/internal/repository"
)

// AdminService proxies flag operations and records audit events.
type AdminService struct {
	flags *repository.FlagsUpstream
	audit *repository.AuditUpstream
}

func NewAdminService(flags *repository.FlagsUpstream, audit *repository.AuditUpstream) *AdminService {
	return &AdminService{flags: flags, audit: audit}
}

func (s *AdminService) List() (status int, body []byte, err error) {
	return s.flags.List()
}

func (s *AdminService) CreateFlag(body []byte) (status int, resp []byte, err error) {
	status, resp, err = s.flags.Create(body)
	if err != nil {
		return status, resp, err
	}
	if status >= 200 && status < 300 {
		var created model.Flag
		if json.Unmarshal(resp, &created) == nil {
			s.logAudit(model.AuditEvent{
				FlagName:  created.Name,
				Action:    "create",
				OldValue:  false,
				NewValue:  created.Enabled,
				ChangedAt: time.Now().UTC(),
			})
		}
	}
	return status, resp, nil
}

func (s *AdminService) ToggleFlag(name string) (status int, resp []byte, err error) {
	oldValue, _ := s.fetchFlagEnabled(name)
	status, resp, err = s.flags.Toggle(name)
	if err != nil {
		return status, resp, err
	}
	if status >= 200 && status < 300 {
		var toggled model.Flag
		if json.Unmarshal(resp, &toggled) == nil {
			s.logAudit(model.AuditEvent{
				FlagName:  name,
				Action:    "toggle",
				OldValue:  oldValue,
				NewValue:  toggled.Enabled,
				ChangedAt: time.Now().UTC(),
			})
		}
	}
	return status, resp, nil
}

func (s *AdminService) DeleteFlag(name string) (status int, resp []byte, err error) {
	oldValue, _ := s.fetchFlagEnabled(name)
	status, resp, err = s.flags.Delete(name)
	if err != nil {
		return status, resp, err
	}
	if status >= 200 && status < 300 {
		s.logAudit(model.AuditEvent{
			FlagName:  name,
			Action:    "delete",
			OldValue:  oldValue,
			NewValue:  false,
			ChangedAt: time.Now().UTC(),
		})
	}
	return status, resp, nil
}

func (s *AdminService) fetchFlagEnabled(name string) (bool, error) {
	status, body, err := s.flags.Get(name)
	if err != nil {
		return false, err
	}
	if status != http.StatusOK {
		return false, nil
	}
	var flag model.Flag
	if err := json.Unmarshal(body, &flag); err != nil {
		return false, err
	}
	return flag.Enabled, nil
}

func (s *AdminService) logAudit(event model.AuditEvent) {
	payload, err := json.Marshal(event)
	if err != nil {
		log.Printf("failed to marshal audit event: %v", err)
		return
	}
	if _, _, err := s.audit.Log(payload); err != nil {
		log.Printf("failed to write audit event: %v", err)
	}
}
