package netrc

import "io/ioutil"
import "fmt"

// Writer encapsulates the concept of writing an netrc file
type Writer struct {
}

// WriterInterface defines the Reader functions
type WriterInterface interface {
	Write(n Interface) error
}

// Writer attempts to read the netrc file
func (w *Writer) Write(n Interface) error {
	path, err := getPath()
	if err != nil {
		return err
	}
	machines := n.GetMachines()
	data := ""
	for _, m := range machines {
		if len(m.HostName) == 0 {
			continue
		}
		data = fmt.Sprintf("%smachine %s\n", data, m.HostName)
		if len(m.DefaultHostName) > 0 {
			data = fmt.Sprintf("%s  default %s\n", data, m.DefaultHostName)
		}
		if len(m.UserName) > 0 {
			data = fmt.Sprintf("%s  login %s\n", data, m.UserName)
		}
		if len(m.Password) > 0 {
			data = fmt.Sprintf("%s  password %s\n", data, m.Password)
		}
		if len(m.AccountPassword) > 0 {
			data = fmt.Sprintf("%s  account %s\n", data, m.AccountPassword)
		}
		if len(m.MacroName) > 0 {
			data = fmt.Sprintf("%s  macdef %s\n", data, m.MacroName)
		}
	}
	if len(data) > 0 {
		return ioutil.WriteFile(path, []byte(data), 0644)
	}
	return nil
}
