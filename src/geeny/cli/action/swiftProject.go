package action

import (
	cli "geeny/cli"
	output "geeny/output"

	random "github.com/Pallinder/go-randomdata"
)

// SwiftProject cli action
func (a *Action) SwiftProject(c *cli.Context) (*cli.Meta, error) {
	if c == nil {
		panic("context is nil")
	}
	name, err := c.GetString(0)
	if err != nil {
		randName := random.SillyName()
		name = &randName
	}

	return a.swiftProject(*name)
}

// - private

func (a *Action) swiftProject(name string) (*cli.Meta, error) {
	// First do geeny swift thing
	meta, err := a.swiftThing(name)
	if err != nil {
		return nil, err
	}
	cID := meta.Items[0].ID
	cName := meta.Items[0].Name
	ttID := meta.Items[1].ID
	ttName := meta.Items[1].Name
	tID := meta.Items[2].ID
	pairingCode := meta.Items[2].Info

	// Create an app, with a made-up name
	cmd, err := a.Tree.CommandForPath([]string{"geeny", "apps", "create"}, 0)
	if err != nil {
		return nil, err
	}
	stage := "dev"
	err = cmd.SetValueForOptionWithFlag(&name, "n")
	if err != nil {
		return nil, err
	}
	err = cmd.SetValueForOptionWithFlag(&stage, "s")
	if err != nil {
		return nil, err
	}
	meta, err = cmd.Exec()
	if err != nil {
		return nil, err
	}
	item, err := meta.ItemFromRawJSON("id", "name", "")
	if err != nil {
		return nil, err
	}
	appID := item.ID
	appName := item.Name

	// Create a client, with a made-up name
	cmd, err = a.Tree.CommandForPath([]string{"geeny", "clients", "create"}, 0)
	if err != nil {
		return nil, err
	}
	secret := "1234"
	clientName := name + "-client"
	redirectURI := "http://test.com"
	scopes := "myscope"
	err = cmd.SetValueForOptionWithFlag(&appID, "id")
	if err != nil {
		return nil, err
	}
	err = cmd.SetValueForOptionWithFlag(&secret, "c")
	if err != nil {
		return nil, err
	}
	err = cmd.SetValueForOptionWithFlag(&clientName, "n")
	if err != nil {
		return nil, err
	}
	err = cmd.SetValueForOptionWithFlag(&redirectURI, "r")
	if err != nil {
		return nil, err
	}
	err = cmd.SetValueForOptionWithFlag(&scopes, "s")
	if err != nil {
		return nil, err
	}
	meta, err = cmd.Exec()
	if err != nil {
		return nil, err
	}
	item, err = meta.ItemFromRawJSON("id", "", "")
	if err != nil {
		return nil, err
	}
	clientID := item.ID

	// Authorize app
	cmd, err = a.Tree.CommandForPath([]string{"geeny", "auths", "create"}, 0)
	if err != nil {
		return nil, err
	}
	scopes = "things/" + tID
	err = cmd.SetValueForOptionWithFlag(&cID, "id")
	if err != nil {
		return nil, err
	}
	err = cmd.SetValueForOptionWithFlag(&scopes, "s")
	if err != nil {
		return nil, err
	}
	meta, err = cmd.Exec()
	if err != nil {
		return nil, err
	}

	// Create a project
	cmd, err = a.Tree.CommandForPath([]string{"geeny", "projects", "create"}, 0)
	if err != nil {
		return nil, err
	}
	projectName := name + "-project"
	err = cmd.SetValueForOptionWithFlag(&appID, "aid")
	if err != nil {
		return nil, err
	}
	err = cmd.SetValueForOptionWithFlag(&projectName, "n")
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
	projID := item.ID
	projName := item.Name

	// Generate project repo
	meta, err = a.generateProject(projID)
	if err != nil {
		return nil, err
	}
	repoName := meta.Items[0].Name

	// Link mediation
	cmd, err = a.Tree.CommandForPath([]string{"geeny", "mediations", "create"}, 0)
	if err != nil {
		return nil, err
	}
	handlerName := "my-first-mediation-handler"
	err = cmd.SetValueForOptionWithFlag(&cID, "id")
	if err != nil {
		return nil, err
	}
	err = cmd.SetValueForOptionWithFlag(&handlerName, "n")
	if err != nil {
		return nil, err
	}
	err = cmd.SetValueForOptionWithFlag(&cID, "cid")
	if err != nil {
		return nil, err
	}
	err = cmd.SetValueForOptionWithFlag(&projID, "pid")
	if err != nil {
		return nil, err
	}
	meta, err = cmd.Exec()
	if err != nil {
		return nil, err
	}
	item, err = meta.ItemFromRawJSON("id", "", "")
	if err != nil {
		return nil, err
	}
	mediationID := item.ID

	// Link pipeline
	cmd, err = a.Tree.CommandForPath([]string{"geeny", "pipelines", "create"}, 0)
	if err != nil {
		return nil, err
	}
	actorName := "my-first-actor"
	definition := "DefaultActor," + actorName
	err = cmd.SetValueForOptionWithFlag(&appID, "id")
	if err != nil {
		return nil, err
	}
	err = cmd.SetValueForOptionWithFlag(&cID, "cid")
	if err != nil {
		return nil, err
	}
	err = cmd.SetValueForOptionWithFlag(&definition, "d") //TODO: cant find this in swagger because description is missing
	if err != nil {
		return nil, err
	}
	err = cmd.SetValueForOptionWithFlag(&actorName, "n")
	if err != nil {
		return nil, err
	}
	meta, err = cmd.Exec()
	if err != nil {
		return nil, err
	}
	item, err = meta.ItemFromRawJSON("", "name", "")
	if err != nil {
		return nil, err
	}
	pipelineName := item.Name

	// Create hook
	cmd, err = a.Tree.CommandForPath([]string{"geeny", "hooks", "create"}, 0)
	if err != nil {
		return nil, err
	}
	payload := "hello world!"
	url := "http://api.geeny.io/demohook"
	err = cmd.SetValueForOptionWithFlag(&appID, "id")
	if err != nil {
		return nil, err
	}
	err = cmd.SetValueForOptionWithFlag(&payload, "p")
	if err != nil {
		return nil, err
	}
	err = cmd.SetValueForOptionWithFlag(&url, "u")
	if err != nil {
		return nil, err
	}
	meta, err = cmd.Exec()
	if err != nil {
		return nil, err
	}
	item, err = meta.ItemFromRawJSON("id", "", "url")
	if err != nil {
		return nil, err
	}
	webhookID := item.ID
	webhookURL := item.ID

	output.Println("Now you can `cd ", repoName+"`", "and run `geeny projects deploy -m \"first\"`")
	output.Println("You can then test it by sending data with `geeny things test -tid", tID, "-cid", cID, "-s`")

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
			&cli.Item{
				ID:   appID,
				Name: appName,
			},
			&cli.Item{
				ID:   clientID,
				Name: clientName,
			},
			&cli.Item{
				ID:   projID,
				Name: projName,
			},
			&cli.Item{
				Name: repoName,
			},
			&cli.Item{
				ID: mediationID,
			},
			&cli.Item{
				Name: pipelineName,
			},
			&cli.Item{
				ID:   webhookID,
				Info: webhookURL,
			},
		},
	}, nil
}
