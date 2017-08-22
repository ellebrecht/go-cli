package action

import (
	cli "geeny/cli"
	output "geeny/output"
	version "geeny/version"
)

// CheckUpdate cli action
func (a *Action) CheckUpdate(c *cli.Context) (*cli.Meta, error) {
	return a.checkUpdate()
}

// - private

func (a *Action) checkUpdate() (*cli.Meta, error) {
	if version.UpdateChecked {
		return nil, nil
	}
	spinner := output.NewSpinner()
	spinner.Text(false, "checking for update")
	spinner.Start()
	defer spinner.Stop(false)
	newVersion, message := version.CheckUpdate(false)
	if newVersion {
		spinner.Text(false, message)
	} else {
		spinner.Text(false, "You already have the latest version "+version.Version)
	}
	return &cli.Meta{
		Info: version.Version,
	}, nil
}
