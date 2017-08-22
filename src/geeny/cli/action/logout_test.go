package action

//TODO
/*
import (
	"errors"
	"testing"

	cli "geeny/cli"
	tu "testing/util"
)

func TestLogoutNoNetrcError(t *testing.T) {
	netrcReader := tu.MockNetrcReader{
		Error: errors.New("read error"),
	}
	action := &Action{
		APIManager:  new(tu.MockAPIManager),
		NetrcReader: netrcReader,
	}
	_, err := action.Logout(new(cli.Context))
	if err != netrcReader.Error {
		t.Fatal("expected error:", netrcReader.Error, "got:", err)
	}
}

func TestLogoutSaveNetrcError(t *testing.T) {
	netrc := &tu.MockNetrc{
		Error: errors.New("save error"),
	}
	netrcReader := tu.MockNetrcReader{
		Netrc: netrc,
	}
	action := &Action{
		APIManager:  new(tu.MockAPIManager),
		NetrcReader: netrcReader,
	}
	_, err := action.Logout(new(cli.Context))
	if err != netrc.Error {
		t.Fatal("expected error:", netrc.Error, "got:", err)
	}
}

func TestLogoutSuccess(t *testing.T) {
	netrcReader := tu.MockNetrcReader{
		Netrc: new(tu.MockNetrc),
	}
	action := &Action{
		APIManager:  new(tu.MockAPIManager),
		NetrcReader: netrcReader,
	}
	_, err := action.Logout(new(cli.Context))
	if err != nil {
		t.Fatal("got error:", err)
	}
}
*/