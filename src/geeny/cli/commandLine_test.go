package cli

import (
	"strconv"
	"testing"
)

func TestSearch(t *testing.T) {
	commandLine := NewCommandLine()
	testCases := []struct {
		args     []string
		flag     string
		expected bool
	}{
		{[]string{"-flag1", "-flag2", "-flag3", "-flag4"}, "flag2", true},
		{[]string{"-flag1", "--flag2", "-flag3", "-flag4"}, "flag2", true},
		{[]string{"-flag1", "-flag2", "-flag3", "-flag4"}, "flag5", false},
		{[]string{"-flag1", "-flag2", "-flag3", "-flag4"}, "bad!", false},
	}
	for _, tc := range testCases {
		res := commandLine.Search(tc.args, tc.flag)
		if res != tc.expected {
			t.Fatal("expected:", strconv.FormatBool(tc.expected), "got:", strconv.FormatBool(res))
		}
	}
}
