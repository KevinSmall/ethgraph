package logr

import (
	"io"
	"log"
	"os"
)

var (
	Trace   *log.Logger // Needs verbose flag. Just about anything
	Info    *log.Logger // Always displayed. Important information
	Warning *log.Logger // Always displayed. Be concerned
	Error   *log.Logger // Displayed with additional location info. Critical problem
)

var isVerbose = false
var target *io.Writer

func init() {
	SetTarget(os.Stdout)
	SetVerbosity(isVerbose)
}

// SetVerbosity allows verbose logs
func SetVerbosity(isVerboseRequested bool) {
	if isVerboseRequested {
		isVerbose = true
	} else {
		isVerbose = false
	}
	refresh()
}

// SetTarget for log output
func SetTarget(newTarget io.Writer) {
	target = &newTarget
	refresh()
}

func refresh() {
	// Trace needs verbose flag to be visible
	if isVerbose {
		Trace = log.New(*target, "", 0)
	} else {
		Trace = log.New(io.Discard, "", 0)
	}

	// Info and Warn always displayed
	Info = log.New(*target, "", 0)
	Warning = log.New(*target, "WARNING: ", 0)

	// Errors always display, and give more detail
	Error = log.New(*target, "ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
