package action

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	model "geeny/api/model"
	"geeny/cli"
	"geeny/output"
)

// DeployProject cli action
func (a *Action) DeployProject(c *cli.Context) (*cli.Meta, error) {
	if c == nil {
		panic("context is nil")
	}
	if c.Count() < 1 {
		return nil, errors.New("expected 1 argument")
	}
	message, err := c.GetStringForFlag("m")
	if err != nil {
		return nil, err
	}

	return a.deployProject(*message)
}

// - private

func (a *Action) deployProject(message string) (*cli.Meta, error) {
	spinner := output.NewSpinner()
	spinner.Text(false, "deploying project, please wait...")
	spinner.Start()
	defer spinner.Stop(false)

	// get github project name
	projectID, err := a.GitHub.GetGeenyConfigValue(".", "projectId")
	if err != nil {
		return nil, err
	}

	root, _ := filepath.Abs(".")
	files, _ := filepath.Glob("./*/*/build.sbt")

	for _, value := range files {
		dir := filepath.Dir(value)
		spinner.Text(false, "Unit testing "+dir+"...")
		os.Chdir(dir)
		out, err := exec.Command("sbt", "test").Output()
		os.Chdir(root)
		if err != nil {
			fmt.Printf("'sbt test' failed in %s:\n %s", dir, out)
			return nil, err
		}
	}

	// get all projects
	cmd, err := a.Tree.CommandForPath([]string{"geeny", "projects", "list"}, 0)
	if err != nil {
		return nil, err
	}

	var meta *cli.Meta
	output.DisableForAction(func() {
		meta, err = cmd.Exec()
	})
	if err != nil {
		return nil, err
	}
	items, err := meta.ItemsFromRawJSON("id", "name", "repositoryName")
	if err != nil {
		return nil, err
	}

	// get github project name
	spinner.Text(false, "getting github config...")
	projectID, err = a.GitHub.GetGeenyConfigValue(".", "projectId")
	if err != nil {
		return nil, err
	}

	// find project that matches github project name
	var projMeta *cli.Item
	for _, item := range items {
		if strings.Compare(item.ID, *projectID) == 0 {
			projMeta = item
			break
		}
	}
	if projMeta == nil {
		return nil, errors.New("there is no longer a project with id: " + *projectID + ". Has it recently been deleted?")
	}

	// get the project's repo information
	spinner.Text(false, "getting repo information...")
	cmd, err = a.Tree.CommandForPath([]string{"geeny", "projects", "get"}, 0)
	if err != nil {
		return nil, err
	}
	err = cmd.SetValueForOptionWithFlag(&projMeta.ID, "id")
	if err != nil {
		return nil, err
	}
	output.DisableForAction(func() {
		meta, err = cmd.Exec()
	})
	if err != nil {
		return nil, err
	}
	repo := &model.Repository{
		Name: "",
	}
	err = meta.UnmarshalRawJSONAtPath(&repo, "repository")
	if err != nil {
		return nil, err
	}
	if repo == nil || len(repo.Name) == 0 {
		return nil, errors.New("the repository '" + projMeta.Info + "' no longer exists")
	} else if !repo.IsValid() {
		return nil, errors.New("the repository was malformed from the api")
	}

	// prepare repo
	spinner.Text(false, "preparing repo...")
	keyFile, err := a.GitHub.SetUpRepoSSH(repo)
	if err != nil {
		return nil, err
	}
	defer a.GitHub.TearDownRepoSSH(repo, keyFile)

	// push repo
	spinner.Text(false, "pushing repo...")
	_ = a.GitHub.Add("*", ".")
	_ = a.GitHub.Commit(message, ".")
	err = a.GitHub.PushRepo(repo)
	if err != nil {
		return nil, err
	}

	spinner.Text(false, "project deployment has started, please check here for updates: https://dashboard.geeny.io/projects/"+projMeta.ID+"/status-board")
	return &cli.Meta{
		Items: []*cli.Item{
			&cli.Item{
				ID:   projMeta.ID,
				Name: projMeta.Info,
			},
		},
	}, nil
}
