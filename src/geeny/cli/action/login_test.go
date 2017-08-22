package action

//TODO
/*
import (
	"errors"
	"testing"

	response "geeny/api/response"
	cli "geeny/cli"
	tu "testing/util"
)

func TestLoginArgsCountError(t *testing.T) {
	action := &Action{
		ConnectManager: new(tu.MockAPIManager),
	}
	context := &cli.Context{
		Args: []*cli.Option{},
	}
	expectedErrMsg := "expected 2 arguments"
	_, err := action.Login(context)
	if err == nil || err.Error() != expectedErrMsg {
		t.Fatal("expected error:", expectedErrMsg, "got:", err)
	}
}

func TestLoginArgsError(t *testing.T) {
	testLoginArgs(t, "test@test.com", "", "password is missing")
	testLoginArgs(t, "", "testPass", "email is missing")
}

func TestLoginNetworkError(t *testing.T) {
	apiManager := &tu.MockAPIManager{
		Error: errors.New("network error"),
	}
	action := &Action{
		ConnectManager: apiManager,
	}
	_, err := action.Login(testLoginValidContext())
	if err != apiManager.Error {
		t.Fatal("expected error:", apiManager.Error, "got:", err)
	}
}

func TestLoginBadLoginError(t *testing.T) {
	apiManager := &tu.MockAPIManager{
		Payloads: []interface{}{
			&response.LoginResponse{
				IDToken:   "",
				TokenType: "Bearer",
			},
		},
	}
	action := &Action{
		ConnectManager: apiManager,
	}
	expectedErrMsg := "bad api token. are your credentials correct?"
	_, err := action.Login(testLoginValidContext())
	if err == nil || err.Error() != expectedErrMsg {
		t.Fatal("expected error:", expectedErrMsg, "got:", err)
	}
}

func TestLoginNoNetrcError(t *testing.T) {
	apiManager := testLoginValidAPIManager()
	netrcReader := tu.MockNetrcReader{
		Error: errors.New("read error"),
	}
	action := &Action{
		ConnectManager:  apiManager,
		NetrcReader: netrcReader,
	}
	_, err := action.Login(testLoginValidContext())
	if err != netrcReader.Error {
		t.Fatal("expected error:", netrcReader.Error, "got:", err)
	}
}

func TestLoginSaveNetrcError(t *testing.T) {
	apiManager := testLoginValidAPIManager()
	netrc := &tu.MockNetrc{
		Error: errors.New("save error"),
	}
	netrcReader := tu.MockNetrcReader{
		Netrc: netrc,
	}
	action := &Action{
		ConnectManager:  apiManager,
		NetrcReader: netrcReader,
	}
	_, err := action.Login(testLoginValidContext())
	if err != netrc.Error {
		t.Fatal("expected error:", netrc.Error, "got:", err)
	}
}

func TestLoginSuccess(t *testing.T) {
	apiManager := testLoginValidAPIManager()
	action := &Action{
		ConnectManager: apiManager,
		NetrcReader: &tu.MockNetrcReader{
			Netrc: new(tu.MockNetrc),
		},
	}
	_, err := action.Login(testLoginValidContext())
	if err != nil {
		t.Fatal("got error:", err)
	}
}

// - private

func testLoginArgs(t *testing.T, email string, pass string, expected string) {
	action := &Action{
		ConnectManager: new(tu.MockAPIManager),
	}
	context := &cli.Context{
		Args: []*cli.Option{
			&cli.Option{Name: "email", Value: &email},
			&cli.Option{Name: "password", Value: &pass},
		},
	}
	_, err := action.Login(context)
	if err == nil || err.Error() != expected {
		t.Fatal("expected error:", expected, "got:", err)
	}
}

func testLoginValidAPIManager() *tu.MockAPIManager {
	return &tu.MockAPIManager{
		Payloads: []interface{}{
			&response.LoginResponse{
				IDToken:   "testToken",
				TokenType: "Bearer",
			},
		},
	}
}

func testLoginValidContext() *cli.Context {
	email := "testEmail@gmail.com"
	pass := "testPassword"
	return &cli.Context{
		Args: []*cli.Option{
			&cli.Option{Value: &email},
			&cli.Option{Value: &pass},
		},
	}
}
*/
