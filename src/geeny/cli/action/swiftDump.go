package action

import (
	"geeny/cli"
	"geeny/output"
)

// SwiftDump cli action
func (a *Action) SwiftDump(c *cli.Context) (*cli.Meta, error) {
	return a.swiftDump()
}

// - private

func (a *Action) swiftDump() (*cli.Meta, error) {
	// content types
	cmd, _ := a.Tree.CommandForPath([]string{"geeny", "content-types", "list"}, 0)
	var cTypesMeta *cli.Meta
	if cmd != nil {
		cTypesMeta, _ = cmd.Exec()
	}

	// thing types
	cmd, _ = a.Tree.CommandForPath([]string{"geeny", "thing-types", "list"}, 0)
	var tTypesMeta *cli.Meta
	if cmd != nil {
		tTypesMeta, _ = cmd.Exec()
	}

	// things
	cmd, _ = a.Tree.CommandForPath([]string{"geeny", "things", "list"}, 0)
	if cmd != nil {
		_, _ = cmd.Exec()
	}

	// projects
	cmd, _ = a.Tree.CommandForPath([]string{"geeny", "projects", "list"}, 0)
	if cmd != nil {
		_, _ = cmd.Exec()
	}

	// apps
	cmd, _ = a.Tree.CommandForPath([]string{"geeny", "apps", "list"}, 0)
	var appsMeta *cli.Meta
	if cmd != nil {
		appsMeta, _ = cmd.Exec()
	}

	// addons
	cmd, _ = a.Tree.CommandForPath([]string{"geeny", "addons", "list"}, 0)
	if cmd != nil {
		_, _ = cmd.Exec()
	}

	// mediations
	if cTypesMeta != nil {
		output.Println("\n---- Mediations ----")
		items, err := cTypesMeta.ItemsFromRawJSON("id", "name", "")
		if err != nil {
			return nil, err
		}
		for _, i := range items {
			output.Println("\ncontent type: " + i.Name)
			cmd, _ := a.Tree.CommandForPath([]string{"geeny", "mediations", "list"}, 0)
			if cmd != nil {
				err = cmd.SetValueForOptionWithFlag(&i.ID, "cid")
				if err != nil {
					return nil, err
				}
				_, _ = cmd.Exec()
			}
		}
	}

	// pipelines
	if appsMeta != nil {
		output.Println("\n---- Pipelines ----")
		items, err := appsMeta.ItemsFromRawJSON("id", "name", "")
		if err != nil {
			return nil, err
		}
		for _, i := range items {
			output.Println("\napp: " + i.Name)
			cmd, _ := a.Tree.CommandForPath([]string{"geeny", "piplines", "list"}, 0)
			if cmd != nil {
				err = cmd.SetValueForOptionWithFlag(&i.ID, "id")
				if err != nil {
					return nil, err
				}
				_, _ = cmd.Exec()
			}
		}
	}

	// firmwares
	if tTypesMeta != nil {
		output.Println("\n---- Firmwares ----")
		items, err := tTypesMeta.ItemsFromRawJSON("id", "name", "")
		if err != nil {
			return nil, err
		}
		version := "0.0.0"
		for _, i := range items {
			output.Println("\nthing type: " + i.Name)
			cmd, _ := a.Tree.CommandForPath([]string{"geeny", "firmwares", "list"}, 0)
			if cmd != nil {
				err = cmd.SetValueForOptionWithFlag(&i.ID, "id")
				if err != nil {
					return nil, err
				}
				err = cmd.SetValueForOptionWithFlag(&version, "v")
				if err != nil {
					return nil, err
				}
				_, _ = cmd.Exec()
			}
		}
	}

	return nil, nil
}
