package model

import "strings"

// Repository represents the data type of a repository
type Repository struct {
	Name         string  `json:"name"`
	State        string  `json:"state"`
	ProjectID    string  `json:"projectId"`
	RepositoryID float32 `json:"repositoryId"`
	PublicKey    string  `json:"publicKey"`
	PrivateKey   string  `json:"privateKey"`
	CreatedAt    string  `json:"createdAt"`
	UpdatedAt    string  `json:"updatedAt"`
	URL          string  `json:"repositoryUrl"`
}

// IsSSH returns true if ssh
func (r *Repository) IsSSH() bool {
	return strings.Contains(r.URL, "ssh")
}

// IsHTTPS returns true if https
func (r *Repository) IsHTTPS() bool {
	return strings.Contains(r.URL, "https")
}

// - ValidationInterface

// IsValid validates the data structure
func (r *Repository) IsValid() bool {
	return (len(r.Name) > 0 &&
		len(r.State) > 0 &&
		len(r.ProjectID) > 0 &&
		r.RepositoryID > 0 &&
		len(r.PublicKey) > 0 &&
		len(r.PrivateKey) > 0 &&
		len(r.CreatedAt) > 0 &&
		len(r.UpdatedAt) > 0 &&
		len(r.URL) > 0)
}
