package action

//TODO
/*
import (
	"errors"
	"os"
	"testing"

	model "geeny/api/model"
	cli "geeny/cli"
	tu "testing/util"
)

func TestGenerateProjectNoTokenError(t *testing.T) {
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
	_, err := action.GenerateProject(testGenerateProjectValidContext())
	if err != netrc.Error {
		t.Fatal("expected error:", netrc.Error, "got:", err)
	}
}

func TestGenerateProjectNetworkError(t *testing.T) {
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
	_, err := action.GenerateProject(testGenerateProjectValidContext())
	if err != apiManager.Error {
		t.Fatal("expected error:", apiManager.Error, "got:", err)
	}
}

func TestGenerateProjectArgsError(t *testing.T) {
	testGenerateProjectArgs(t, "", "projectID is missing")
}

func TestGenerateProjectArgCountError(t *testing.T) {
	action := &Action{
		APIManager:  new(tu.MockAPIManager),
		NetrcReader: new(tu.MockNetrcReader),
	}
	context := &cli.Context{
		Args: []*cli.Option{},
	}
	_, err := action.GenerateProject(context)
	expected := "expected 1 argument"
	if err == nil || err.Error() != expected {
		t.Fatal("expected error:", expected, "got:", err)
	}
}

func TestGenerateProjectGithubError(t *testing.T) {
	apiManager := testGenerateProjectValidAPIManager()
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
	_, err := action.GenerateProject(testGenerateProjectValidContext())
	if err != gitHub.Error {
		t.Fatal("expected error:", gitHub.Error, "got:", err)
	}
}

func TestGenerateProjectSuccess(t *testing.T) {
	testGenerateProjectCreateDummyFiles()
	defer testGenerateProjectDeleteDummyFiles()
	apiManager := testGenerateProjectValidAPIManager()
	netrc := &tu.MockNetrc{
		Pwd: "test",
	}
	netrcReader := tu.MockNetrcReader{
		Netrc: netrc,
	}
	gitHub := tu.MockGitHub{}
	action := &Action{
		APIManager:  apiManager,
		NetrcReader: netrcReader,
		GitHub:      gitHub,
	}
	_, err := action.GenerateProject(testGenerateProjectValidContext())
	if err != nil {
		t.Fatal("got error:", err)
	}
}

// - private

func testGenerateProjectArgs(t *testing.T, projectID string, expected string) {
	apiManager := testGenerateProjectValidAPIManager()
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
			&cli.Option{Name: "projectID", Value: &projectID},
		},
	}
	_, err := action.GenerateProject(context)
	if err == nil || err.Error() != expected {
		t.Fatal("expected error:", expected, "got:", err)
	}
}

func testGenerateProjectValidAPIManager() *tu.MockAPIManager {
	project := testGenerateProjectVaidProject()
	return &tu.MockAPIManager{
		Payloads: []interface{}{
			project,
		},
	}
}

func testGenerateProjectVaidProject() *model.Project {
	return &model.Project{
		Repository: model.Repository{
			Name:         "test",
			State:        "test",
			ProjectID:    "test",
			RepositoryID: 1,
			PublicKey:    "test",
			PrivateKey:   "test",
			CreatedAt:    "test",
			UpdatedAt:    "test",
			URL:          "http://test",
		},
	}
}

func testGenerateProjectValidContext() *cli.Context {
	projectID := "testProjectID"
	return &cli.Context{
		Args: []*cli.Option{
			&cli.Option{Value: &projectID},
		},
	}
}

func testGenerateProjectCreateDummyFiles() {
	os.Chdir("/tmp")
	_ = os.MkdirAll("/tmp/geeny-io-project-sample", 0777)
	_ = os.MkdirAll(testGenerateProjectVaidProject().Repository.Name, 0777)
}

func testGenerateProjectDeleteDummyFiles() {
	os.Chdir("/tmp")
	_ = os.RemoveAll("/tmp/geeny-io-project-sample")
	_ = os.RemoveAll(testGenerateProjectVaidProject().Repository.Name)
}
*/
