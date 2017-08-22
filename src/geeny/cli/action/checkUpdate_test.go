package action

import "testing"

func TestCheckUpdateSuccess(t *testing.T) {
	action := &Action{}
	_, err := action.CheckUpdate(nil)
	if err != nil {
		t.Fatal("got:", err)
	}
}
