package cli

import (
	"errors"
	"testing"
)

var errCmd = errors.New("test")
var action = func(context *Context) (*Meta, error) {
	return nil, nil
}
var fail = func(context *Context) (*Meta, error) {
	return nil, errCmd
}
var tree = &Command{
	Name:   "1",
	Action: action,
	Options: []*Option{
		&Option{
			Name: "Test",
			Flag: "t",
		},
	},
	Commands: []*Command{
		&Command{
			Name:   "2",
			Action: action,
			Options: []*Option{
				&Option{
					Name: "Test",
					Flag: "t",
				},
			},
			Commands: []*Command{
				&Command{
					Name:        "3",
					Summary:     "",
					Action:      fail,
					Options: []*Option{
						&Option{
							Name: "Test",
							Flag: "t",
						},
					},
				},
			},
		},
	},
}

func TestRootCommandUsageError(t *testing.T) {
	parser := newParser(true)
	_, err := parser.parseCommand(tree, []string{"1"}, 0)
	if err == nil {
		t.Fatal("expected error")
	}
	expected := errUsage
	if err.errCode != expected {
		t.Fatalf("expected error code: %d but got: %d", expected, err.errCode)
	}
}

func TestNestedCommandUsageError(t *testing.T) {
	parser := newParser(true)
	_, err := parser.parseCommand(tree, []string{"1", "2"}, 0)
	if err == nil {
		t.Fatal("expected error")
	}
	expected := errUsage
	if err.errCode != expected {
		t.Fatalf("expected error code: %d but got: %d", expected, err.errCode)
	}
}

func TestCommandError(t *testing.T) {
	parser := newParser(true)
	_, err := parser.parseCommand(tree, []string{"1", "2", "3", "-t", "test"}, 0)
	if err == nil {
		t.Fatal("expected error")
	}
	expectedCode := errCommandAction
	if err.errCode != expectedCode {
		t.Fatalf("expected error code: %d but got: %d", expectedCode, err.errCode)
	}
	expectedErr := errCmd
	if err.internalErr != expectedErr {
		t.Fatalf("expected error code: %v but got: %d", expectedErr, err.errCode)
	}
}

func TestRootCommandExecutes(t *testing.T) {
	parser := newParser(true)
	_, err := parser.parseCommand(tree, []string{"1", "-t", "test"}, 0)
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
}

func TestNestedCommandExecutes(t *testing.T) {
	parser := newParser(true)
	_, err := parser.parseCommand(tree, []string{"1", "2", "-t", "test"}, 0)
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
}

func TestHelpFlag(t *testing.T) {
	parser := newParser(true)
	_, err := parser.parseCommand(tree, []string{"1", "2", "3", "-h"}, 0)
	expected := errUsage
	if err.errCode != expected {
		t.Fatalf("expected error code: %d but got: %d", expected, err.errCode)
	}
}

func TestBashCompletionFlag(t *testing.T) {
	parser := newParser(true)
	_, err := parser.parseCommand(tree, []string{"1", "2", "--generate-bash-completion"}, 0)
	expected := errBashCompletion
	if err.errCode != expected {
		t.Fatalf("expected error code: %d but got: %d", expected, err.errCode)
	}
}
