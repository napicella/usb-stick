package interactive

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
)

func PromptForPassword() (string, error) {
	fmt.Print("Enter Password: ")
	bytePassword, e := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if e != nil {
		return "", e
	}

	return string(bytePassword), nil
}
