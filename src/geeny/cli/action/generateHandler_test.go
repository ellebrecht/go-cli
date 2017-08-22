package action

//TODO
/*
import (
	cli "geeny/cli"
	"os"
	"testing"

	tu "testing/util"
)

func TestGenerateHandlerArgsError(t *testing.T) {
	testGenerateHandlerArgs(t, "", "name is missing")
}

func TestGenerateHandlerArgCountError(t *testing.T) {
	action := &Action{}
	context := &cli.Context{
		Args: []*cli.Option{},
	}
	_, err := action.GenerateHandler(context)
	expected := "expected 1 argument"
	if err == nil || err.Error() != expected {
		t.Fatal("expected error:", expected, "got:", err)
	}
}

func TestGenerateHandlerNotInRepoError(t *testing.T) {
	gitHub := tu.MockGitHub{
		IsRepoValue: false,
	}
	action := &Action{
		GitHub: gitHub,
	}
	_, err := action.GenerateHandler(testGenerateHandlerValidContext())
	expected := "you are not in a git repo"
	if err == nil || err.Error() != expected {
		t.Fatal("expected error:", expected, "got:", err)
	}
}

func TestGenerateHandlerExistsError(t *testing.T) {
	testGenerateHandlerCreateExistingDirectory("testName")
	defer testGenerateHandlerDeleteDummyFiles()
	gitHub := tu.MockGitHub{
		IsRepoValue: true,
	}
	action := &Action{
		GitHub: gitHub,
	}
	_, err := action.GenerateHandler(testGenerateHandlerValidContext())
	expected := "a mediation handler already exists with this name, please choose another"
	if err == nil || err.Error() != expected {
		t.Fatal("expected error:", expected, "got:", err)
	}
}

func TestGenerateHandlerSuccess(t *testing.T) {
	testGenerateHandlerCreateDummyFiles()
	defer testGenerateHandlerDeleteDummyFiles()
	gitHub := tu.MockGitHub{
		IsRepoValue: true,
	}
	action := &Action{
		GitHub: gitHub,
	}
	_, err := action.GenerateHandler(testGenerateHandlerValidContext())
	if err != nil {
		t.Fatal("got error:", err)
	}
}

// - private

func testGenerateHandlerArgs(t *testing.T, name string, expected string) {
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
	_, err := action.GenerateHandler(context)
	if err == nil || err.Error() != expected {
		t.Fatal("expected error:", expected, "got:", err)
	}
}

func testGenerateHandlerValidContext() *cli.Context {
	name := "testName"
	return &cli.Context{
		Args: []*cli.Option{
			&cli.Option{Value: &name},
		},
	}
}

func testGenerateHandlerCreateDummyFiles() {
	os.Chdir("/tmp")
	_ = os.MkdirAll("mediation", 0777)
	_ = os.MkdirAll("geeny-io-mediation-sample", 0777)
}

func testGenerateHandlerDeleteDummyFiles() {
	os.Chdir("/tmp")
	_ = os.RemoveAll("mediation")
	_ = os.RemoveAll("geeny-io-mediation-sample")
}

func testGenerateHandlerCreateExistingDirectory(name string) {
	os.Chdir("/tmp")
	_ = os.MkdirAll("mediation/"+name, 0777)
}
*/
