package odf

import (
	"flag"
	"log"
	"os"
)

var logger = log.New(os.Stderr, "[odf] ", log.LstdFlags)

var debug = flag.Bool("odf-debug", false, "Allow debugging output to os.Stderr")

func printf(format string, v ...interface{}) {
	if debug == nil || !*debug {
		return
	}
	logger.Printf(format, v...)
}
