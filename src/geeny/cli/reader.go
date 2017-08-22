package cli

import (
	"bufio"
	"os"
	"strings"
	"syscall"

	"geeny/output"

	"golang.org/x/crypto/ssh/terminal"
)

// ReadSecure reads data from the command line whilst obfuscating output
func ReadSecure() string {
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	for err != nil {
		output.Println("Error:", err)
		return ReadSecure()
	}
	return strings.TrimSpace(string(bytePassword))
}

// Read reads data from the command line (not secure)
func Read() string {
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	for err != nil {
		output.Println("Error:", err)
		return Read()
	}
	return strings.TrimSpace(text)
}
