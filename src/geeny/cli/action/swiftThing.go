package action

import (
	"geeny/cli"

	random "github.com/Pallinder/go-randomdata"
)

// SwiftThing cli action
func (a *Action) SwiftThing(c *cli.Context) (*cli.Meta, error) {
	if c == nil {
		panic("context is nil")
	}
	name, err := c.GetString(0)
	if err != nil {
		return a.swiftThing(random.SillyName())
	}
	return a.swiftThing(*name)
}

// - private

func (a *Action) swiftThing(name string) (*cli.Meta, error) {
	// Create a content-type using a made-up name
	cmd, err := a.Tree.CommandForPath([]string{"geeny", "content-types", "create"}, 0)
	if err != nil {
		return nil, err
	}
	desc := "created with geeny swift"
	err = cmd.SetValueForOptionWithFlag(&desc, "d")
	if err != nil {
		return nil, err
	}
	err = cmd.SetValueForOptionWithFlag(&name, "n")
	if err != nil {
		return nil, err
	}
	meta, err := cmd.Exec()
	if err != nil {
		return nil, err
	}
	item, err := meta.ItemFromRawJSON("id", "name", "")
	if err != nil {
		return nil, err
	}
	cID := item.ID
	cName := item.Name

	// Create a thing-type with made-up name, using the content-type
	cmd, err = a.Tree.CommandForPath([]string{"geeny", "thing-types", "create"}, 0)
	if err != nil {
		return nil, err
	}
	err = cmd.SetValueForOptionWithFlag(&cID, "ids")
	if err != nil {
		return nil, err
	}
	err = cmd.SetValueForOptionWithFlag(&name, "n")
	if err != nil {
		return nil, err
	}
	meta, err = cmd.Exec()
	if err != nil {
		return nil, err
	}
	item, err = meta.ItemFromRawJSON("id", "name", "")
	if err != nil {
		return nil, err
	}
	ttID := item.ID
	ttName := item.Name

	// Creates a thing using the thing-type
	meta, err = a.createThing(ttID, true)
	if err != nil {
		return nil, err
	}
	tID := meta.Items[0].ID
	pairingCode := meta.Items[0].Info

	// Pair the thing to the user account
	cmd, err = a.Tree.CommandForPath([]string{"geeny", "things", "pair"}, 0)
	if err != nil {
		return nil, err
	}
	err = cmd.SetValueForOptionWithFlag(&pairingCode, "p")
	if err != nil {
		return nil, err
	}
	err = cmd.SetValueForOptionWithFlag(&name, "n")
	if err != nil {
		return nil, err
	}
	_, err = cmd.Exec()
	if err != nil {
		return nil, err
	}

	return &cli.Meta{
		Items: []*cli.Item{
			&cli.Item{
				ID:   cID,
				Name: cName,
			},
			&cli.Item{
				ID:   ttID,
				Name: ttName,
			},
			&cli.Item{
				ID:   tID,
				Info: pairingCode,
			},
		},
	}, nil
}
