package action

import (
	"errors"
	"os"
	"time"

	model "geeny/api/model"
	cli "geeny/cli"
	github "geeny/github"
	io "geeny/io"
	output "geeny/output"
)

// GenerateProject cli action
func (a *Action) GenerateProject(c *cli.Context) (*cli.Meta, error) {
	if c == nil {
		panic("context is nil")
	}
	if c.Count() < 1 {
		return nil, errors.New("expected 1 argument")
	}
	projectID, err := c.GetStringForFlag("id")
	if err != nil {
		return nil, err
	}

	return a.generateProject(*projectID)
}

// - private

func (a *Action) generateProject(projectID string) (*cli.Meta, error) {
	spinner := output.NewSpinner()
	spinner.Text(false, "waiting for the project's repo to be created...")
	spinner.Start()
	defer spinner.Stop(false)

	// wait until the project has a repo
	repo, err := a.getRepoWhenCreated(projectID, spinner)
	if err != nil {
		return nil, err
	}

	// init empty repo
	spinner.Text(true, "initialising local repo, please wait...")
	err = os.MkdirAll(repo.Name+"/.git", 0777)
	if err != nil {
		return nil, err
	}
	err = a.GitHub.InitRepo(repo.Name + "/.git")
	if err != nil {
		return nil, err
	}

	// set config values
	err = a.GitHub.SetConfigValue(repo.Name, "remote.origin.url", repo.URL)
	if err != nil {
		return nil, err
	}
	err = a.GitHub.SetGeenyConfigValue(repo.Name, "projectId", projectID)
	if err != nil {
		return nil, err
	}

	// clone template repo
	spinner.Text(true, "cloning template project, please wait...")
	templateProjName := "geeny-io-project-sample"
	err = a.GitHub.CloneRepo(&model.Repository{
		Name: templateProjName,
		URL:  github.HTTPS("quodio", templateProjName),
	})
	if err != nil {
		return nil, err
	}
	// delete git history
	err = os.RemoveAll(templateProjName + "/.git")
	if err != nil {
		return nil, err
	}
	// copy template
	err = io.CopyDirContents(templateProjName, repo.Name)
	if err != nil {
		return nil, err
	}
	// remove template
	err = os.RemoveAll(templateProjName)
	if err != nil {
		return nil, err
	}
	// rename remote
	_ = a.GitHub.RenameRemote("origin", "geeny", repo.Name) //TODO: github issue - this will fail, but actually succeed, so we ignore the error until github fix the issue
	// prepare repo for push
	err = a.GitHub.Add("*", repo.Name)
	err = a.GitHub.Commit("commit from: geeny generate project", repo.Name)

	spinner.Text(false,
		"setup a new project template named", repo.Name+".", "to push this (and your future changes), use 'projects deploy'\n\n",
		"You must create a new repo in your GitHub account, then run the following commands:\n",
		"  $ cd", repo.Name+"\n",
		"  $ git remote add origin git@github.com:{your organization}/{repo name}.git\n",
		"  $ git push -u origin master\n")
	return &cli.Meta{
		Items: []*cli.Item{
			&cli.Item{
				ID:   projectID,
				Name: repo.Name,
			},
		},
	}, nil
}

func (a *Action) getRepoWhenCreated(projectID string, spinner *output.Spinner) (*model.Repository, error) {
	cmd, err := a.Tree.CommandForPath([]string{"geeny", "projects", "get"}, 0)
	if err != nil {
		return nil, err
	}
	var meta *cli.Meta
	err = cmd.SetValueForOptionWithFlag(&projectID, "id")
	if err != nil {
		return nil, err
	}
	output.DisableForAction(func() {
		meta, err = cmd.Exec()
	})
	if err != nil {
		return nil, err
	}
	repo := model.Repository{}
	err = meta.UnmarshalRawJSONAtPath(&repo, "repository")
	if err != nil {
		return nil, err
	}

	// Project.Repository won't be created until the repo is created in github
	// if it isn't there, keep asking until it's there
	// todo: this will be improved in the api, perhaps a separate call projects/<id>/repo/status
	if len(repo.Name) == 0 {
		time.Sleep(5 * time.Second)
		spinner.Text(false, "still waiting for repo to be initialized...")
		return a.getRepoWhenCreated(projectID, spinner)
	} else if repo.IsValid() == false {
		return nil, errors.New("malformed repository from api")
	}
	return &repo, nil
}
