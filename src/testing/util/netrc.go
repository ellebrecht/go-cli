package util

import netrc "geeny/netrc"

// - Mock Netrc

type MockNetrc struct {
	Pwd      string
	Error    error
	Machines []*netrc.Machine
}

func (n MockNetrc) Remove(account string) {}
func (n MockNetrc) Password(account string) (*string, error) {
	return &n.Pwd, n.Error
}
func (n MockNetrc) Add(account string, email string, password string) {}
func (n MockNetrc) GetMachines() []*netrc.Machine {
	return n.Machines
}

// - Mock NetrcReader

type MockNetrcReader struct {
	Netrc *MockNetrc
	Error error
}

func (n MockNetrcReader) Read() (netrc.Interface, error) {
	return n.Netrc, n.Error
}
