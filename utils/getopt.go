package utils

import (
	"os"
	"flag"
	"strings"
	"github.com/gwtony/gapi/variable"
)

// EOF means reach the end
const EOF = -1

// OptErr Option error
var OptErr = 1
// OptInd Option index
var OptInd = 1

// OptOpt Current option
var OptOpt uint8
// OptArg Current option args
var OptArg string

var sp = 1

// Getopt gets option
func Getopt(opts string) int {
	var c uint8
	var cp int
	argv := os.Args
	argc := len(argv)

	if sp == 1 {
		if OptInd >= argc ||
			(len(argv[0]) > 0 && argv[OptInd][0] != '-') ||
			len(argv[0]) == 1 {
			return EOF
		} else if argv[OptInd] == "--" {
			OptInd++
			return EOF
		}
	}
	c = argv[OptInd][sp]
	OptOpt = c
	cp = strings.Index(opts, string(c))
	if c == ':' || cp == -1 {
		if OptErr != 0 {
			println(": illegal option --", string(c))
		}
		sp++
		if len(argv[OptInd]) == sp {
			OptInd++
			sp = 1
		}
		return '?'
	}
	cp++
	if cp < len(opts) && opts[cp] == ':' {
		if len(argv[OptInd]) > sp+1 {
			OptArg = argv[OptInd][sp+1 : sp+2]
			OptInd++
		} else {
			OptInd++
			if OptInd >= argc {
				if OptErr != 0 {
					println(": option requires an argument --", string(c))
				}
				sp = 1
				return '?'
			}
			OptArg = argv[OptInd]
			OptInd++
		}
		sp = 1
	} else {
		sp++
		if len(argv[OptInd]) == sp {
			sp = 1
			OptInd++
		}
		OptArg = ""
	}
	return int(c)
}

var (
	ConfigFile = flag.String("f", "", "Set config file to macedon")
	Version    = flag.Bool("v", false, "Show version")
)

func ParseOption() int {
	flag.Parse()
	if *Version {
		println("Version is", variable.VERSION)
		return -1
	}

	return 0
}
