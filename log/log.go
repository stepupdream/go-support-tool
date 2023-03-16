//go:build !linux

package log

import (
	"io"
	"log"
	"os"

	"github.com/stepupdream/golang-support-tool/util"
)

func Setting(filename string, isDebug bool) {
	// Open file for write/read logging. (if not, generate one)
	logfile, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	// Generate a Writer for both normal and file output.
	multiLogFile := io.MultiWriter(os.Stdout, logfile)

	// Log output settings (display date and time)
	// Adding log.Llongfile will also output the log output points.
	if isDebug {
		log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	} else {
		log.SetFlags(log.Ldate | log.Ltime)
	}

	// Specify log output destination.
	log.SetOutput(multiLogFile)
}

func Fatal(messages ...any) {
	for _, message := range messages {
		log.Printf("\x1b[31m%s\x1b[0m\n", message)
	}
	util.KeyWait("Press any key:")
	os.Exit(1)
}
