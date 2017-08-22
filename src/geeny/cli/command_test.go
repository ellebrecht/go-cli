package cli

import (
	"fmt"
	"strings"
	"testing"
)

func TestCommandForPathSuccess(t *testing.T) {
	testTree := testTree()
	testCases := []struct {
		input    []string
		index    int
		expected string
	}{
		{[]string{"1", "2a", "3a", "4"}, 0, "TestCommandForPath1"},
		{[]string{"1", "2a", "3b", "4b"}, 0, "TestCommandForPath2"},
		{[]string{"1", "2b", "3", "4b"}, 0, "TestCommandForPath3"},
	}
	for _, tc := range testCases {
		test, err := testTree.CommandForPath(tc.input, tc.index)
		if err != nil {
			t.Fatal("got err:", err)
		}
		if strings.Compare(test.Summary, tc.expected) != 0 {
			t.Fatal("expected:", tc.expected, "got:", test.Summary)
		}
	}
}

func TestCommandForPathFail(t *testing.T) {
	testTree := testTree()
	_, err := testTree.CommandForPath([]string{"1", "bad!"}, 0)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSubCommandForNameSuccess(t *testing.T) {
	testTree := testTree()
	cmd := testTree.SubCommandForName("2a")
	expected := "TestSubCommandForName"
	if strings.Compare(cmd.Summary, expected) != 0 {
		t.Fatal("expected:", expected, "got:", cmd.Summary)
	}
}

func TestSubCommandForNameFail(t *testing.T) {
	testTree := testTree()
	cmd := testTree.SubCommandForName("bad!")
	if cmd != nil {
		t.Fatal("expected nil")
	}
}

func TestSubCommandForPathSuccess(t *testing.T) {
	testTree := testTree()
	cmd, res := testTree.SubCommandForPath([]string{"2b", "3", "4c"})
	if !res {
		t.Fatal("expected true")
	}
	expected := "TestSubCommandForPath"
	if strings.Compare(cmd.Summary, expected) != 0 {
		t.Fatal("expected:", expected, "got:", cmd.Summary)
	}
}

func TestSubCommandForPathFail(t *testing.T) {
	testTree := testTree()
	cmd, res := testTree.SubCommandForPath([]string{"bad!"})
	if res {
		t.Fatal("expected false")
	}
	if cmd != nil {
		t.Fatal("expected nil")
	}
}

func TestOptionForFlagSuccess(t *testing.T) {
	testTree := testTree()
	opt, err := testTree.OptionForFlag("f1")
	if err != nil {
		t.Fatal("got err:", err)
	}
	expected := "TestOptionForFlag"
	if strings.Compare(opt.Description, expected) != 0 {
		t.Fatal("expected:", expected, "got:", opt.Description)
	}
}

func TestOptionForFlagFail(t *testing.T) {
	testTree := testTree()
	_, err := testTree.OptionForFlag("bad!")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSetValueForOptionWithFlagSuccess(t *testing.T) {
	testTree := testTree()
	expected := "test"
	err := testTree.SetValueForOptionWithFlag(&expected, "f1")
	if err != nil {
		t.Fatal("got err:", err)
	}
	opt := testTree.Options[0]
	if opt.Value != &expected {
		t.Fatal("expected:", expected, "got:", opt.Value)
	}
}

func TestSetValueForOptionWithFlagFail(t *testing.T) {
	testTree := testTree()
	err := testTree.SetValueForOptionWithFlag("test", "f1")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestExec(t *testing.T) {
	val1 := "testVal1"
	val2 := "testVal2"
	cmd := &Command{
		Options: []*Option{
			&Option{
				Flag:  "f1",
				Value: &val1,
			},
			&Option{
				Flag:  "f2",
				Value: &val2,
			},
		},
		Action: func(c *Context) (*Meta, error) {
			if len(c.Args) != 2 {
				t.Fatal("expected 2 args")
			}
			testCases := []struct {
				input    string
				expected string
			}{
				{*c.Args[0].Value.(*string), val1},
				{*c.Args[1].Value.(*string), val2},
			}
			for _, tc := range testCases {
				if strings.Compare(tc.input, tc.expected) != 0 {
					t.Fatal("expected:", val1, "got:", tc.input)
				}
			}
			return nil, nil
		},
	}
	_, _ = cmd.Exec()
}

func TestSetupParents(t *testing.T) {
	testTree := testTree()
	testTree.SetupParents()
	testParentOfCommand(testTree, t)
	if testTree.Commands[0].Hidden == false {
		t.Fatal("expected true")
	}
}

func TestMergeSuccess(t *testing.T) {
	tree1 := testTree()
	tree2 := &Command{
		Name: "1",
		Commands: []*Command{
			&Command{
				Name: "2a",
				Commands: []*Command{
					&Command{
						Name:    "3c",
						Summary: "TestMerge1",
					},
				},
			},
			&Command{
				Name: "2b",
				Commands: []*Command{
					&Command{
						Name:    "3b",
						Summary: "TestMerge2",
					},
				},
			},
		},
	}

	fmt.Printf("old:\n%s\n\n", tree1.TreeString(0))
	for _, sub := range tree2.Commands {
		err := tree1.Merge(sub)
		if err != nil {
			t.Fatal("got err:", err)
		}
	}
	fmt.Printf("new:\n%s", tree1.TreeString(0))

	expected1 := "TestMerge1"
	expected2 := "TestMerge2"
	testCases := []struct {
		input    []string
		index    int
		expected *string
	}{
		// old tree
		{[]string{"1", "2a", "3a", "4"}, 0, nil},
		{[]string{"1", "2a", "3b", "4a"}, 0, nil},
		{[]string{"1", "2a", "3b", "4b"}, 0, nil},
		{[]string{"1", "2b", "3", "4a"}, 0, nil},
		{[]string{"1", "2b", "3", "4b"}, 0, nil},
		{[]string{"1", "2b", "3", "4c"}, 0, nil},

		// new elements
		{[]string{"1", "2a", "3c"}, 0, &expected1},
		{[]string{"1", "2b", "3b"}, 0, &expected2},
	}
	for _, tc := range testCases {
		test, err := tree1.CommandForPath(tc.input, tc.index)
		if err != nil {
			t.Fatal("got err:", err)
		}
		if tc.expected != nil && strings.Compare(test.Summary, *tc.expected) != 0 {
			t.Fatal("expected:", tc.expected, "got:", test.Summary)
		}
	}
}

// - private

func testParentOfCommand(cmd *Command, t *testing.T) {
	for _, sub := range cmd.Commands {
		if sub.Parent != cmd {
			t.Fatal("expected:", cmd, "got:", sub.Parent)
		}
		testParentOfCommand(sub, t)
	}
}

func testTree() *Command {
	return &Command{
		Name: "1",
		Options: []*Option{
			&Option{
				Flag:        "f1",
				Description: "TestOptionForFlag",
			},
			&Option{
				Flag: "f2",
			},
		},
		Commands: []*Command{
			&Command{
				Name:    "2a",
				Summary: "TestSubCommandForName",
				Commands: []*Command{
					&Command{
						Name:   "3a",
						Hidden: true,
						Commands: []*Command{
							&Command{
								Name:     "4",
								Summary:  "TestCommandForPath1",
								Commands: []*Command{},
							},
						},
					},
					&Command{
						Name:   "3b",
						Hidden: true,
						Commands: []*Command{
							&Command{
								Name:     "4a",
								Commands: []*Command{},
							},
							&Command{
								Name:     "4b",
								Summary:  "TestCommandForPath2",
								Commands: []*Command{},
							},
						},
					},
				},
			},
			&Command{
				Name: "2b",
				Commands: []*Command{
					&Command{
						Name: "3",
						Commands: []*Command{
							&Command{
								Name:     "4a",
								Commands: []*Command{},
							},
							&Command{
								Name:     "4b",
								Summary:  "TestCommandForPath3",
								Commands: []*Command{},
							},
							&Command{
								Name:     "4c",
								Summary:  "TestSubCommandForPath",
								Commands: []*Command{},
							},
						},
					},
				},
			},
		},
	}
}
