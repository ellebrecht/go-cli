package netrc

import (
	"errors"
	"strings"

	log "geeny/log"
)

/*
   deals with ~/.netrc file
   example entry:
   machine api.geeny.com
     login roger@health2works.com
     password 8da4a77a-d734-4110-ae24-5cb1b1f44a69`
*/

// Netrc encapsultes the netrc file
type Netrc struct {
	machines []*Machine
}

// Interface defines the Netrc functions
type Interface interface {
	Remove(account string)
	Add(account string, email string, password string)
	Password(account string) (*string, error)
	GetMachines() []*Machine
}

// Instance returns sinleton instance
func Instance() (Interface, error) {
	reader := &Reader{}
	err := reader.Read()
	if err != nil {
		return nil, err
	}
	machines := reader.Parse()
	return &Netrc{
		machines: machines,
	}, nil
}

// Remove removes an account from the netrc file
func (n *Netrc) Remove(account string) {
	log.Warnf("Removing login from netrc for %v", account)
	i, _ := n.machine(account)
	if i >= 0 {
		n.machines = append(n.machines[:i], n.machines[i+1:]...)
	}
}

// Add adds an account from the netrc file
func (n *Netrc) Add(account string, email string, password string) {
	_, m := n.machine(account)
	if m != nil {
		log.Warnf("Adding login to netrc for %v", account)
		m.HostName = account
		m.UserName = email
		m.Password = password
		return
	}
	log.Warnf("Updating login in netrc for %v", account)
	n.machines = append(n.machines, &Machine{
		HostName: account,
		UserName: email,
		Password: password,
	})
}

// Password returns sensitive data (as string) from the netrc file (i.e. token)
func (n *Netrc) Password(account string) (*string, error) {
	_, m := n.machine(account)
	if m == nil {
		log.Warnf("No login found in netrc for %v", account)
		return nil, errors.New("Please login first using \"geeny login\".")
	}
	return &m.Password, nil
}

// GetMachines returns machines
func (n *Netrc) GetMachines() []*Machine {
	return n.machines
}

// - private

func (n *Netrc) machine(name string) (int, *Machine) {
	for i, m := range n.machines {
		if strings.Compare(m.HostName, name) == 0 {
			return i, n.machines[i]
		}
	}
	return -1, nil
}
