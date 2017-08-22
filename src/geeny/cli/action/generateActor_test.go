package action

//TODO
/*
import (
	"os"
	"testing"

	cli "geeny/cli"
	tu "testing/util"
)

func TestGenerateActorArgsError(t *testing.T) {
	testGenerateActorArgs(t, "", "name is missing")
}

func TestGenerateActorArgCountError(t *testing.T) {
	action := &Action{}
	context := &cli.Context{
		Args: []*cli.Option{},
	}
	_, err := action.GenerateActor(context)
	expected := "expected 1 argument"
	if err == nil || err.Error() != expected {
		t.Fatal("expected error:", expected, "got:", err)
	}
}

func TestGenerateActorNotInRepoError(t *testing.T) {
	gitHub := tu.MockGitHub{
		IsRepoValue: false,
	}
	action := &Action{
		GitHub: gitHub,
	}
	_, err := action.GenerateActor(testGenerateActorValidContext())
	expected := "you are not in a git repo"
	if err == nil || err.Error() != expected {
		t.Fatal("expected error:", expected, "got:", err)
	}
}

func TestGenerateActorExistsError(t *testing.T) {
	testGenerateActorCreateExistingDirectory("testName")
	defer testGenerateActorDeleteDummyFiles()
	gitHub := tu.MockGitHub{
		IsRepoValue: true,
	}
	action := &Action{
		GitHub: gitHub,
	}
	_, err := action.GenerateActor(testGenerateActorValidContext())
	expected := "an actor already exists with this name, please choose another"
	if err == nil || err.Error() != expected {
		t.Fatal("expected error:", expected, "got:", err)
	}
}

func TestGenerateActorSuccess(t *testing.T) {
	testGenerateActorCreateDummyFiles()
	defer testGenerateActorDeleteDummyFiles()
	gitHub := tu.MockGitHub{
		IsRepoValue: true,
	}
	action := &Action{
		GitHub: gitHub,
	}
	_, err := action.GenerateActor(testGenerateActorValidContext())
	if err != nil {
		t.Fatal("got error:", err)
	}
}

// - private

func testGenerateActorArgs(t *testing.T, name string, expected string) {
	netrc := &tu.MockNetrc{
		Pwd: "test",
	}
	netrcReader := tu.MockNetrcReader{
		Netrc: netrc,
	}
	action := &Action{
		NetrcReader: netrcReader,
	}
	context := &cli.Context{
		Args: []*cli.Option{
			&cli.Option{Name: "name", Value: &name},
		},
	}
	_, err := action.GenerateActor(context)
	if err == nil || err.Error() != expected {
		t.Fatal("expected error:", expected, "got:", err)
	}
}

func testGenerateActorValidContext() *cli.Context {
	name := "testName"
	return &cli.Context{
		Args: []*cli.Option{
			&cli.Option{Value: &name},
		},
	}
}

func testGenerateActorCreateDummyFiles() {
	os.Chdir("/tmp")
	_ = os.MkdirAll("pipeline", 0777)
	_ = os.MkdirAll("SampleActor", 0777)
}

func testGenerateActorDeleteDummyFiles() {
	os.Chdir("/tmp")
	_ = os.RemoveAll("pipeline")
	_ = os.RemoveAll("SampleActor")
}

func testGenerateActorCreateExistingDirectory(name string) {
	os.Chdir("/tmp")
	_ = os.MkdirAll("pipeline/"+name, 0777)
}
*/
