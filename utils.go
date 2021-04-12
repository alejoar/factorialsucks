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

const URL_SIGN_IN = "https://factorialhr.com/users/sign_in"
const URL_CLOCK_IN = "https://app.factorialhr.com/attendance/clock-in"
const API_URL = "https://api.factorialhr.com"

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var WEEKEND_DAYS = [...]string{
	"Saturday",
	"Sunday",
	"sábado",
	"domingo",
	"dissabte",
	"diumenge",
	"zaterdag",
	"zondag",
	"sabato",
	"domenica",
	"lördag",
	"söndag",
	"Samstag",
	"Sonntag",
	"samedi",
	"dimanche",
	"Sábado",
	"Domingo",
}

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
