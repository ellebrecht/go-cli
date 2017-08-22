package action

import (
	"errors"
	"io/ioutil"
	"os"
	"os/user"

	"geeny/api/model"
	"geeny/cli"
	"geeny/output"
)

// CreateThing cli action
func (a *Action) CreateThing(c *cli.Context) (*cli.Meta, error) {
	if c == nil {
		panic("context is nil")
	}
	if c.Count() < 2 {
		return nil, errors.New("expected 2 arguments")
	}
	thingTypeID, err := c.GetStringForFlag("id")
	if err != nil {
		return nil, err
	}
	storeSecrets, err := c.GetBoolForFlag("s")
	if err != nil {
		return nil, err
	}

	return a.createThing(*thingTypeID, *storeSecrets)
}

// - private

func (a *Action) createThing(ttID string, storeSecrets bool) (*cli.Meta, error) {
	// run command
	cmd, err := a.Tree.CommandForPath([]string{"geeny", "things", "createbase"}, 0)
	if err != nil {
		return nil, err
	}
	err = cmd.SetValueForOptionWithFlag(&storeSecrets, "s")
	if err != nil {
		return nil, err
	}
	err = cmd.SetValueForOptionWithFlag(&ttID, "id")
	if err != nil {
		return nil, err
	}
	meta, err := cmd.Exec()
	if err != nil {
		return nil, err
	}

	// get info from thing and generate certificate
	item, err := meta.ItemFromRawJSON("id", "", "pairingCode")
	tID := item.ID
	pairingCode := item.Info
	thing := model.Thing{}
	err = meta.UnmarshalRawJSON(&thing)
	if err != nil {
		return nil, err
	}
	pem := thing.Attributes.IOT.Certificate.CertificatePem
	key := thing.Attributes.IOT.Certificate.KeyPair.PrivateKey

	// make ~/.geeny/.../ dirs if they dont exist
	usr, err := user.Current()
	if err != nil {
		return nil, errors.New("Unable to determine your home directory")
	}
	path := (usr.HomeDir + string(os.PathSeparator) +
		".geeny" + string(os.PathSeparator) +
		"things" + string(os.PathSeparator) +
		tID)
	err = os.MkdirAll(path, 0700)
	if err != nil {
		return nil, err
	}

	// write certs
	pemPath := path + string(os.PathSeparator) + "cert.pem"
	keyPath := path + string(os.PathSeparator) + "cert.key"
	err = ioutil.WriteFile(pemPath, []byte(pem), 0644)
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(keyPath, []byte(key), 0644)
	if err != nil {
		return nil, err
	}

	output.Println("Your thing's certificate has been saved at: " + path)

	return &cli.Meta{
		Items: []*cli.Item{
			&cli.Item{
				ID:   tID,
				Info: pairingCode,
			},
		},
	}, nil
}
