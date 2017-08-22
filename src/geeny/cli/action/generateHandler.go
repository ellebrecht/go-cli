package action

import (
	"errors"
	"os"

	model "geeny/api/model"
	cli "geeny/cli"
	github "geeny/github"
	io "geeny/io"
	output "geeny/output"
)

// GenerateHandler cli action
func (a *Action) GenerateHandler(c *cli.Context) (*cli.Meta, error) {
	if c == nil {
		panic("context is nil")
	}
	if c.Count() < 1 {
		return nil, errors.New("expected 1 argument")
	}
	name, err := c.GetString(0)
	if err != nil {
		return nil, err
	}

	return a.generateHandler(*name)
}

// - private

func (a *Action) generateHandler(name string) (*cli.Meta, error) {
	// check if git repo
	if !a.GitHub.IsRepo(".") {
		return nil, errors.New("you are not in a git repo")
	}

	// check if geeny project
	_, err := a.GitHub.GetGeenyConfigValue(".", "projectId")
	if err != nil {
		return nil, errors.New("you are not in a geeny project, or projectId is missing from .git/config")
	}

	// check if it's already a handler?
	path := "mediation/" + name
	_, err = os.Stat(path)
	if err == nil {
		return nil, errors.New("a mediation handler already exists with this name, please choose another")
	}

	spinner := output.NewSpinner()
	spinner.Text(false, "generating handler, please wait...")
	spinner.Start()
	defer spinner.Stop(false)

	// clone template repo
	repoName := "geeny-io-mediation-sample"
	err = a.GitHub.CloneRepo(&model.Repository{
		Name: repoName,
		URL:  github.HTTPS("quodio", repoName),
	})
	if err != nil {
		return nil, err
	}

	// move directory
	err = io.MoveDir(repoName, path)
	if err != nil {
		return nil, err
	}

	spinner.Text(false, path+" was generated. use 'projects deploy' to push changes")
	return &cli.Meta{
		Items: []*cli.Item{
			&cli.Item{
				Name: name,
			},
		},
	}, nil
}
