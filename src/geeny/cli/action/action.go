package action

import (
	cli "geeny/cli"
	gitub "geeny/github"
	mqtt "geeny/mqtt"
	netrc "geeny/netrc"
)

// Action contains the context for an action
type Action struct {
	APIURL         string
	ConnectURL     string
	TimeoutSeconds float32
	Tree           *cli.Command
	Netrc          netrc.Interface
	NetrcWriter    netrc.WriterInterface
	GitHub         gitub.GitHubInterface
	MQTT           mqtt.MQTTInterface
}

// NewAction creates a new cli action
func NewAction(apiUrl string, connectUrl string) *Action {
	n, err := netrc.Instance()
	if err != nil {
		panic(err)
	}
	return &Action{
		APIURL:         apiUrl,
		ConnectURL:     connectUrl,
		TimeoutSeconds: 30.0,
		Netrc:          n,
		NetrcWriter:    &netrc.Writer{},
		GitHub:         &gitub.GitHub{},
		MQTT:           &mqtt.MQTT{},
	}
}
