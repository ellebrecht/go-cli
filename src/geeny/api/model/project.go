package model

// Project represents the data type of a project
type Project struct {
	ID             string     `json:"id"`
	IdentityID     string     `json:"identityId"`
	Name           string     `json:"name"`
	State          string     `json:"state"`
	CreatedAt      string     `json:"createdAt"`
	ModifiedAt     string     `json:"modifiedAt"`
	RepositoryName string     `json:"repositoryName"`
	Repository     Repository `json:"repository"`
}

// HasRepo returns `true` if the a github repo has been created for the project
func (p *Project) HasRepo() bool {
	return len(p.Repository.Name) > 0
}

// - ValidationInterface

// IsValid validates the data structure
func (p *Project) IsValid() bool {
	return (len(p.ID) > 0 &&
		len(p.IdentityID) > 0 &&
		len(p.Name) > 0 &&
		len(p.CreatedAt) > 0 &&
		len(p.ModifiedAt) > 0 &&
		len(p.RepositoryName) > 0)
}
