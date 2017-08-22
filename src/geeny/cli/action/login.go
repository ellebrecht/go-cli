package action

import (
	"errors"

	"geeny/cli"
	"geeny/output"
)

// Login cli action
func (a *Action) Login(c *cli.Context) (*cli.Meta, error) {
	if c == nil {
		panic("context is nil")
	}
	if c.Count() < 2 {
		return nil, errors.New("expected 2 arguments")
	}
	email, err := c.GetString(0)
	if err != nil {
		return nil, err
	}
	password, err := c.GetString(1)
	if err != nil {
		return nil, err
	}

	return a.login(*email, *password)
}

// - private

func (a *Action) login(email string, password string) (*cli.Meta, error) {
	spinner := output.NewSpinner()
	spinner.Text(false, "logging in, please wait...")
	spinner.Start()
	defer spinner.Stop(false)

	cmd, err := a.Tree.CommandForPath([]string{"geeny", "loginbase"}, 0)
	if err != nil {
		return nil, err
	}
	err = cmd.SetValueForOptionWithFlag(&email, "e")
	if err != nil {
		return nil, err
	}
	err = cmd.SetValueForOptionWithFlag(&password, "p")
	if err != nil {
		return nil, err
	}
	var meta *cli.Meta
	output.DisableForAction(func() {
		meta, err = cmd.Exec()
	})
	if err != nil {
		return nil, err
	}
	item, err := meta.ItemFromRawJSON("", "", "id_token")
	if err != nil {
		return nil, err
	}
	idToken := item.Info
	if len(idToken) == 0 {
		return nil, errors.New("bad api token. are your credentials correct?")
	}

	a.Netrc.Remove(a.APIURL)
	a.Netrc.Add(a.APIURL, email, idToken)

	a.Netrc.Remove(a.ConnectURL)
	a.Netrc.Add(a.ConnectURL, email, idToken)

	err = a.NetrcWriter.Write(a.Netrc)
	if err != nil {
		return nil, err
	}

	spinner.Text(false, "logged in as", email)
	return &cli.Meta{
		Items: []*cli.Item{
			&cli.Item{
				ID:   idToken,
				Name: email,
			},
		},
	}, nil
}
