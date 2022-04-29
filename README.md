# factorialsucks

```
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

## Install

```
brew tap alejoar/tap
brew install factorialsucks
````

## Update

```
brew update
brew upgrade factorialsucks
````
