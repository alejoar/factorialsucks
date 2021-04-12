package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"syscall"

	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh/terminal"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func readCredentials(c *cli.Context) (string, string) {
	email := c.String("email")
	if email == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Email: ")
		email, _ = reader.ReadString('\n')
		email = strings.TrimSuffix(email, "\n")
	}
	if !emailRegex.MatchString(email) {
		log.Fatalln("Email not valid")
	}

	fmt.Print("Password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	password := string(bytePassword)
	if password == "" {
		log.Fatalln("\nNo password provided")
	}
	return email, password
}
