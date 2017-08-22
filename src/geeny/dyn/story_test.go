package dyn

import (
	"path"
	"runtime"
	"testing"

	"geeny/cli"
	"strings"
	tu "testing/util"
)

func TestHappyPath(t *testing.T) { // https://www.youtube.com/watch?v=JmjOkeS9Azo
	// load swagger
	_, b, _, _ := runtime.Caller(0)
	projDir := path.Dir(b)
	plugin := NewPlugin("localhost:1234", "file://"+projDir+"/story_test_swagger.json")
	tree := &cli.Command{
		Name:        "geeny",
		NonCategory: true,
		Commands:    []*cli.Command{},
	}
	tree.SetupParents()
	err := plugin.Init(tree)
	if err != nil {
		t.Fatal("got err:", err)
	}

	// mock network
	jsonBody := `{
		"id": "test"
	}`
	plugin.transport = tu.MockTransport{
		Result: &Response{
			Body: jsonBody,
		},
	}

	// get test command
	//fmt.Println(tree.TreeString(0))
	cmd, err := tree.CommandForPath([]string{"geeny", "apps", "create"}, 0)
	if err != nil {
		t.Fatal("got err:", err)
	}

	// create opts
	opt, err := cmd.OptionForFlag("n")
	if err != nil {
		t.Fatal("got err:", err)
	}
	name := "testName"
	opt.Value = &name

	opt, err = cmd.OptionForFlag("s")
	if err != nil {
		t.Fatal("got err:", err)
	}
	stage := "testStage"
	opt.Value = &stage

	// run with context
	context := &cli.Context{
		Command: cmd,
		Args:    cmd.Options,
	}
	meta, err := cmd.Action(context)
	if err != nil {
		t.Fatal("got err:", err)
	}
	if strings.Compare(jsonBody, meta.RawJSON) != 0 {
		t.Fatal("expected:", jsonBody, "got:", meta.RawJSON)
	}
}
