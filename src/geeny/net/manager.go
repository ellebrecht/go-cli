package net

import (
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	log "geeny/log"
)

// Manager is a thin layer over golang's http.Client
type Manager struct {
	client *http.Client
}

// NewManager creates a new Manager instance
func NewManager(timeout time.Duration) *Manager {
	return &Manager{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// PerformRequest executes a request, and returns a response
func (n *Manager) PerformRequest(req *http.Request) (int, []byte, error) {
	resp, err := n.client.Do(req)
	if resp == nil {
		return 0, nil, errors.New("no response from the api. is there an active internet connection? If this error doesn't go away, please contact support")
	}
	defer resp.Body.Close()
	if err != nil {
		return resp.StatusCode, nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, err
	}
	log.Tracef("Response: status=%v body=%v", resp.StatusCode, string(body))
	return resp.StatusCode, body, nil
}
