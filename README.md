# factorialsucks

```shell
❯ factorialsucks -h

NAME:
   factorialsucks - FactorialHR auto clock in for the whole month from the command line

USAGE:
   factorialsucks [options]

GLOBAL OPTIONS:
   --email value, -e value        you factorial email address
   --year YYYY, -y YYYY           clock-in year YYYY (default: current year)
   --month MM, -m MM              clock-in month MM (default: current month)
   --clock-in HH:MM, --ci HH:MM   clock-in time HH:MM (default: "10:00")
   --clock-out HH:MM, --co HH:MM  clock-in time HH:MM (default: "18:00")
   --today, -t                    clock in for today only (default: false)
   --until-today, --ut            clock in only until today (default: false)
   --dry-run, --dr                do a dry run without actually clocking in (default: false)
   --reset-month, --rm            delete all shifts for the given month (default: false)
   --help, -h                     show help (default: false)
```

[![asciicast](https://asciinema.org/a/1wj0X77lfeHqYWZKp2YY86Xux.svg)](https://asciinema.org/a/1wj0X77lfeHqYWZKp2YY86Xux)

## Build

Using docker/podman, check `build` target in [Makefile] file.

Native build:

```bash
go build -o factorialsucks .
```

## Install

```bash
brew tap alejoar/tap
brew install factorialsucks
```

## Update

```bash
brew update
brew upgrade factorialsucks
```

## Usage

Running `make help` to see the help, options and default values.

If you execute `make run EMAIL=<your@email.com>` this will run the program with
the given email address and the default values.

If you want to run without `Make`. Check `today_splitshift` target in [Makefile]
file as a base command.

For example, to submit the current month:

```bash
docker run \
   -it \
   --rm factorialsucks \
   --email your@email.com \
   --clock-in 7:00 \
   --clock-out 15:00
```

[Makefile]: ./Makefile
