package action

//TODO
/*
import (
	"errors"
	"testing"

	cli "geeny/cli"
	tu "testing/util"
)

func TestStreamLogsNoTokenError(t *testing.T) {
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
	_, err := action.StreamLogs(testStreamLogsValidContext())
	if err != netrc.Error {
		t.Fatal("expected error:", netrc.Error, "got:", err)
	}
}

func TestStreamLogsNetworkError(t *testing.T) {
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
	_, err := action.StreamLogs(testStreamLogsValidContext())
	if err != apiManager.Error {
		t.Fatal("expected error:", apiManager.Error, "got:", err)
	}
}

func TestStreamLogsArgsError(t *testing.T) {
	testStreamLogsArgs(t, "", "", "thingIDs is missing")
}

func TestStreamLogsArgCountError(t *testing.T) {
	action := &Action{
		APIManager:  new(tu.MockAPIManager),
		NetrcReader: new(tu.MockNetrcReader),
	}
	context := &cli.Context{
		Args: []*cli.Option{},
	}
	_, err := action.StreamLogs(context)
	expected := "expected 1 argument"
	if err == nil || err.Error() != expected {
		t.Fatal("expected error:", expected, "got:", err)
	}
}

func TestStreamLogsSuccess(t *testing.T) {
	// TODO: how to test? - it's a long running process
}

// - private

func testStreamLogsArgs(t *testing.T, thingIDs string, serviceIDs, expected string) {
	apiManager := new(tu.MockAPIManager)
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
			&cli.Option{Name: "thingIDs", Value: &thingIDs},
			&cli.Option{Name: "serviceIDs", Value: &serviceIDs},
		},
	}
	_, err := action.StreamLogs(context)
	if err == nil || err.Error() != expected {
		t.Fatal("expected error:", expected, "got:", err)
	}
}

func testStreamLogsValidContext() *cli.Context {
	thingIDs := "testThingID,testThingID"
	serviceIDs := "testServiceID,testServiceID"
	return &cli.Context{
		Args: []*cli.Option{
			&cli.Option{Value: &thingIDs},
			&cli.Option{Value: &serviceIDs},
		},
	}
}
*/
