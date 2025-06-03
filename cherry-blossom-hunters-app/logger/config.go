package logger

import (
	"fmt"
	"log"
	"os"

	"github.com/comail/colog"
)

const (
	Warn  = "WARN"
	Error = "ERROR"
	Debug = "DEBUG"
)

func SetUp() {
	colog.Register()
	colog.SetOutput(os.Stdout)
	colog.ParseFields(true)
	colog.SetFlags(log.Ldate | log.Lshortfile)
}

func Logging(format string, args ...interface{}) {
	level := "info"
	message := format

	// ログレベルを指定する場合は、最初の引数として渡す (例: Logging("warn", "User ID: %v", userID))
	if len(args) > 0 {
		if s, ok := args[0].(string); ok && (s == "warn" || s == "error" || s == "debug" || s == "info") {
			level = s
			if len(args) > 1 {
				message = fmt.Sprintf(format, args[1:]...)
			}
		} else {
			message = fmt.Sprintf(format, args...)
		}
	}

	var prefix string
	switch level {
	case "warn":
		prefix = "warn:"
	case "error":
		prefix = "error:"
	case "debug":
		prefix = "debug:"
	default:
		prefix = "info:"
	}

	log.Printf("%s %s\n", prefix, message)
}
