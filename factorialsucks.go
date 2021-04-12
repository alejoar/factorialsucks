package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/alejoar/factorialsucks/factorial"
	"github.com/urfave/cli/v2"
)

func main() {
	log.SetFlags(0)
	app := &cli.App{
		Name:            "factorialsucks",
		Usage:           "FactorialHR auto clock in for the whole month from the command line",
		Version:         "2.0a",
		Compiled:        time.Now(),
		UsageText:       "factorialsucks [options]",
		HideHelpCommand: true,
		HideVersion:     true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "email",
				Aliases: []string{"e"},
				Usage:   "you factorial email address",
			},
			&cli.Int64Flag{
				Name:        "year",
				Aliases:     []string{"y"},
				Usage:       "clock-in year `YYYY`",
				DefaultText: "current year",
			},
			&cli.Int64Flag{
				Name:        "month",
				Aliases:     []string{"m"},
				Usage:       "clock-in month `MM`",
				DefaultText: "current month",
			},
			&cli.StringFlag{
				Name:    "clock-in",
				Aliases: []string{"ci"},
				Usage:   "clock-in time `HH:MM`",
				Value:   "10:00",
			},
			&cli.StringFlag{
				Name:    "clock-out",
				Aliases: []string{"co"},
				Usage:   "clock-in time `HH:MM`",
				Value:   "18:00",
			},
			&cli.BoolFlag{
				Name:    "dry-run",
				Aliases: []string{"dr"},
				Usage:   "do a dry run without actually clocking in",
			},
		},
		Action: factorialsucks,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func factorialsucks(c *cli.Context) error {
	email, password := readCredentials(c)

	// debug
	fmt.Println()
	fmt.Println("email", email)
	fmt.Println("password", password)
	fmt.Println("year", c.Int64("year"))
	fmt.Println("month", c.Int64("month"))
	fmt.Println("clock-in", c.String("clock-in"))
	fmt.Println("clock-out", c.String("clock-out"))
	fmt.Println("dry-run", c.Bool("dry-run"))
	// /debug

	client := factorial.NewFactorialClient(email, password)
	resp, _ := client.Get("https://api.factorialhr.com/attendance/shifts?year=2021&month=3&employee_id=81638")
	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(data))

	return nil
}
