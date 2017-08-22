package action

import (
	cli "geeny/cli"
	output "geeny/output"
	version "geeny/version"
)

// Version cli action
func (a *Action) Version(c *cli.Context) (*cli.Meta, error) {
	return a.version()
}

// - private

func (a *Action) version() (*cli.Meta, error) {
	output.Println(version.Version)
	return nil, nil
}
