package action

//TODO
/*
import (
	"errors"
	"fmt"
	"testing"

	model "geeny/api/model"
	cli "geeny/cli"
	tu "testing/util"
)

func TestDeployProjectNoTokenError(t *testing.T) {
	netrc := &tu.MockNetrc{
		Error: errors.New("read error"),
	}
	netrcReader := tu.MockNetrcReader{
		Netrc: netrc,
	}
	action := &Action{
		APIManager:  new(tu.MockAPIManager),
		NetrcReader: netrcReader,
	}
	_, err := action.DeployProject(testDeployProjectValidContext())
	if err != netrc.Error {
		t.Fatal("expected error:", netrc.Error, "got:", err)
	}
}

func TestDeployProjectNetworkError(t *testing.T) {
	apiManager := &tu.MockAPIManager{
		Error: errors.New("network error"),
	}
	netrc := &tu.MockNetrc{
		Pwd: "test",
	}
	netrcReader := tu.MockNetrcReader{
		Netrc: netrc,
	}
	action := &Action{
		APIManager:  apiManager,
		NetrcReader: netrcReader,
	}
	_, err := action.DeployProject(testDeployProjectValidContext())
	if err != apiManager.Error {
		t.Fatal("expected error:", apiManager.Error, "got:", err)
	}
}

func TestDeployProjectNameNotFoundError(t *testing.T) {
	project := model.Project{
		Repository: model.Repository{
			Name: "test",
		},
	}
	apiManager := &tu.MockAPIManager{
		Payloads: []interface{}{
			&[]model.Project{
				project,
			},
		},
	}
	netrc := &tu.MockNetrc{
		Pwd: "test",
	}
	netrcReader := tu.MockNetrcReader{
		Netrc: netrc,
	}
	gitHub := tu.MockGitHub{
		Value: "notTest",
	}
	action := &Action{
		APIManager:  apiManager,
		NetrcReader: netrcReader,
		GitHub:      gitHub,
	}
	_, err := action.DeployProject(testDeployProjectValidContext())
	expected := "there is no longer a project with id: " + gitHub.Value + ". Has it recently been deleted?"
	if err == nil || err.Error() != expected {
		t.Fatal("expected error:", expected, "got:", err)
	}
}

func TestDeployProjectNoLongerHostedError(t *testing.T) {
	project := model.Project{
		ID:             "test",
		RepositoryName: "testRepoName",
		Repository:     model.Repository{},
	}
	apiManager := &tu.MockAPIManager{
		Payloads: []interface{}{
			&[]model.Project{
				project,
			},
			&project,
		},
	}
	netrc := &tu.MockNetrc{
		Pwd: "test",
	}
	netrcReader := tu.MockNetrcReader{
		Netrc: netrc,
	}
	gitHub := tu.MockGitHub{
		Value: project.ID,
	}
	action := &Action{
		APIManager:  apiManager,
		NetrcReader: netrcReader,
		GitHub:      gitHub,
	}
	_, err := action.DeployProject(testDeployProjectValidContext())
	expected := "the repository 'testRepoName' no longer exists"
	output.Println(err.Error())
	if err == nil || err.Error() != expected {
		t.Fatal("expected error:", expected, "got:", err)
	}
}

func TestDeployGitHubError(t *testing.T) {
	apiManager := testDeployProjectValidAPIManager()
	netrc := &tu.MockNetrc{
		Pwd: "test",
	}
	netrcReader := tu.MockNetrcReader{
		Netrc: netrc,
	}
	gitHub := tu.MockGitHub{
		Error: errors.New("github error"),
	}
	action := &Action{
		APIManager:  apiManager,
		NetrcReader: netrcReader,
		GitHub:      gitHub,
	}
	_, err := action.DeployProject(testDeployProjectValidContext())
	if err != gitHub.Error {
		t.Fatal("expected error:", gitHub.Error, "got:", err)
	}
}

func TestDeployProjectArgsError(t *testing.T) {
	testDeployProjectArgs(t, "", "message is missing")
}

func TestDeployProjectSuccess(t *testing.T) {
	project := testDeployProjectValidProject()
	apiManager := testDeployProjectValidAPIManager()
	netrc := &tu.MockNetrc{
		Pwd: "test",
	}
	netrcReader := tu.MockNetrcReader{
		Netrc: netrc,
	}
	gitHub := tu.MockGitHub{
		Value: project.ID,
	}
	action := &Action{
		APIManager:  apiManager,
		NetrcReader: netrcReader,
		GitHub:      gitHub,
	}

	_, err := action.DeployProject(testDeployProjectValidContext())
	if err != nil {
		t.Fatal("got error:", err)
	}
}

// - private

func testDeployProjectArgs(t *testing.T, message string, expected string) {
	apiManager := testDeployProjectValidAPIManager()
	netrc := &tu.MockNetrc{
		Pwd: "test",
	}
	netrcReader := tu.MockNetrcReader{
		Netrc: netrc,
	}
	action := &Action{
		APIManager:  apiManager,
		NetrcReader: netrcReader,
	}
	context := &cli.Context{
		Args: []*cli.Option{
			&cli.Option{Name: "message", Value: &message},
		},
	}
	_, err := action.DeployProject(context)
	if err == nil || err.Error() != expected {
		t.Fatal("expected error:", expected, "got:", err)
	}
}

func testDeployProjectValidProject() model.Project {
	return model.Project{
		ID: "test",
		Repository: model.Repository{
			Name:         "test",
			State:        "test",
			ProjectID:    "test",
			RepositoryID: 1,
			PublicKey:    "test",
			PrivateKey:   "test",
			CreatedAt:    "test",
			UpdatedAt:    "test",
			URL:          "https://test",
		},
	}
}

func testDeployProjectValidAPIManager() *tu.MockAPIManager {
	project := testDeployProjectValidProject()
	return &tu.MockAPIManager{
		Payloads: []interface{}{
			&[]model.Project{
				project,
			},
			&project,
		},
	}
}

func testDeployProjectValidContext() *cli.Context {
	message := "testMessage"
	return &cli.Context{
		Args: []*cli.Option{
			&cli.Option{Value: &message},
		},
	}
}
*/
