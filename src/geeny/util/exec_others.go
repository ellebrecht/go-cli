// +build !windows

package util

import (
	log "geeny/log"
	"os/exec"
)

func ExecuteCommand(name string, arg ...string) ([]byte, error) {
	cmd := exec.Command(name, arg...)
	log.Trace("Command exec", *cmd)
	out, err := cmd.Output()
	log.Trace("Command output", string(out))
	if err != nil {
		log.Error("Shell command failed with error: " + err.Error())
	}
	return out, err
}
