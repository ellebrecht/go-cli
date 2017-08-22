package netrc

import (
	"errors"
	"io/ioutil"
	"runtime"
	"strings"

	log "geeny/log"
)

// Reader encapsulates the concept of reading an netrc file
type Reader struct {
	data []byte
}

// ReaderInterface defines the Reader functions
type ReaderInterface interface {
	Read() (Interface, error)
	Parse() (Interface, error)
}

// Read attempts to read the netrc file
func (nr *Reader) Read() error {
	filePath, err := getPath()
	if err != nil {
		return err
	}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Trace("%NETRC env var missing")
		if runtime.GOOS == "windows" {
			return errors.New(err.Error() + ". You might be missing the %NETRC environment variable. try: set NETRC=" + filePath)
		}
		return errors.New(err.Error() + ". You might be missing the missing $NETRC environment variable. try: export NETRC=" + filePath)
	}
	nr.data = data
	return nil
}

// Parse attempts to create Machine objects
func (nr *Reader) Parse() []*Machine {
	data := string(nr.data)
	machineData := strings.Split(data, "machine ")
	machines := []*Machine{}

	for _, m := range machineData[1:] {
		machine := &Machine{}
		machineDetails := strings.Split(m, "\n")
		machine.HostName = strings.TrimSpace(machineDetails[0])

		for _, md := range machineDetails[1:] {
			itemData := strings.Split(md, "default ")
			if len(itemData) == 2 {
				machine.DefaultHostName = strings.TrimSpace(itemData[1])
				continue
			}
			itemData = strings.Split(md, "login ")
			if len(itemData) == 2 {
				machine.UserName = strings.TrimSpace(itemData[1])
				continue
			}
			itemData = strings.Split(md, "password ")
			if len(itemData) == 2 {
				machine.Password = strings.TrimSpace(itemData[1])
				continue
			}
			itemData = strings.Split(md, "account ")
			if len(itemData) == 2 {
				machine.AccountPassword = strings.TrimSpace(itemData[1])
				continue
			}
			itemData = strings.Split(md, "macdef ")
			if len(itemData) == 2 {
				machine.MacroName = strings.TrimSpace(itemData[1])
				continue
			}
		}
		machines = append(machines, machine)
		log.Debug("---- parsed netrc entry ----")
		log.Debug("HostName", machine.HostName)
		log.Debug("DefaultHostName", machine.DefaultHostName)
		log.Debug("UserName", machine.UserName)
		log.Debug("Password", machine.Password)
		log.Debug("AccountPassword", machine.AccountPassword)
		log.Debug("MacroName", machine.MacroName, "\n")
	}
	return machines
}
