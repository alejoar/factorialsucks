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

For native build, you can use the following command:

```bash
go build -o factorialsucks .
```

If your prefer use [Makefile], you can use the following command:

```bash
make build
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

For example, you can use the following command to clock in for the whole month:

```bash
./factorialsucks \
   --email "<your@email.com>" \
   --month 03 \
   --clock-in 10:00 \
   --clock-out 18:00
```

If your prefer use [Makefile], you can use the following command:

```bash
make today_continuous
```

## Docker

### Build with Docker

Using docker/podman, check `build` target in [Makefile] file.

### Use with Docker

There some targets in the [Makefile] to use this program with Docker.

Running `make help` to see the program help, options and default values.

If you execute `make <TARGET> EMAIL=<your@email.com>` this will run the program with the given email address and the default values.

For example, to submit today as split shift with Docker, you can run:

```bash
make docker_today_splitshift
```

[Makefile]: ./Makefile
