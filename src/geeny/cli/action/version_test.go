package action

import "testing"

func TestVersionSuccess(t *testing.T) {
	action := &Action{}
	_, err := action.Version(nil)
	if err != nil {
		t.Fatal("got:", err)
	}
}
