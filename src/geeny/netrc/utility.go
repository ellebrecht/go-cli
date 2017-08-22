package netrc

import (
	"os"
	"os/user"

	log "geeny/log"
)

var getPath = _getPath // @note stubbable for testing

func _getPath() (string, error) {
	path := os.Getenv("NETRC")
	if len(path) == 0 {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		log.Trace("%NETRC env var missing")
		return usr.HomeDir + string(os.PathSeparator) + ".netrc", nil
	}
	log.Trace("%NETRC env var found")
	return path, nil
}
