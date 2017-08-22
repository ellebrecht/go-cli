package plugin

import "geeny/cli"

type Plugin interface {
	Init(root *cli.Command) error
	Close() error
}
