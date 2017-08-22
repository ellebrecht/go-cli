package action

import (
	cli "geeny/cli"
	output "geeny/output"
)

// Logout cli action
func (a *Action) Logout(c *cli.Context) (*cli.Meta, error) {
	return a.logout()
}

// - private

func (a *Action) logout() (*cli.Meta, error) {
	a.Netrc.Remove(a.APIURL)
	a.Netrc.Remove(a.ConnectURL)
	err := a.NetrcWriter.Write(a.Netrc)
	if err != nil {
		return nil, err
	}
	output.Println("logged out")

	return &cli.Meta{
		Items: []*cli.Item{
			&cli.Item{
				Info: a.APIURL,
			},
			&cli.Item{
				Info: a.ConnectURL,
			},
		},
	}, nil
}
