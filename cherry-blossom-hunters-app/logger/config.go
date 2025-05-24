package logger

import (
	"log"
	"os"

	"github.com/comail/colog"
)


const (
	Warn  = "WARN"
	Error = "ERROR"
)

func SetUp() {
	colog.Register()
	colog.SetOutput(os.Stdout)
	colog.ParseFields(true)
	colog.SetFlags(log.Ldate | log.Lshortfile)
}

func Logging(message string, level ...string) {
	lv := "info"

	if len(level) > 0 {
		lv = level[0]
	}

	var prefix string
	switch lv {
	case Warn:
		prefix = "warn:"
	case Error:
		prefix = "error:"
	default:
		prefix = "info:"
	}

	log.Printf("%s %s\n", prefix, message)
}
