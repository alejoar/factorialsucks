package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
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
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
