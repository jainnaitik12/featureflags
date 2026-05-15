package repository

import (
	"fmt"
	"net/http"
	"net/url"

	"featureflags/admin-service/internal/client"
)

// FlagsUpstream calls flags-service over HTTP.
type FlagsUpstream struct {
	baseURL string
}

func NewFlagsUpstream(baseURL string) *FlagsUpstream {
	return &FlagsUpstream{baseURL: baseURL}
}

func (u *FlagsUpstream) List() (int, []byte, error) {
	return client.Do(http.MethodGet, u.baseURL+"/flags", nil)
}

func (u *FlagsUpstream) Create(body []byte) (int, []byte, error) {
	return client.Do(http.MethodPost, u.baseURL+"/flags", body)
}

func (u *FlagsUpstream) Toggle(name string) (int, []byte, error) {
	enc := url.PathEscape(name)
	return client.Do(http.MethodPatch, fmt.Sprintf("%s/flags/%s/toggle", u.baseURL, enc), nil)
}

func (u *FlagsUpstream) Delete(name string) (int, []byte, error) {
	enc := url.PathEscape(name)
	return client.Do(http.MethodDelete, fmt.Sprintf("%s/flags/%s", u.baseURL, enc), nil)
}

func (u *FlagsUpstream) Get(name string) (int, []byte, error) {
	enc := url.PathEscape(name)
	return client.Do(http.MethodGet, fmt.Sprintf("%s/flags/%s", u.baseURL, enc), nil)
}
