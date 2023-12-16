package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Logger = zerolog.New(zerolog.ConsoleWriter{
	Out:        os.Stdout,
	TimeFormat: time.RFC3339,
}).Level(zerolog.TraceLevel).With().Timestamp().Logger()

func Info() *zerolog.Event {
	return Logger.Info()
}

func Warn() *zerolog.Event {
	return Logger.Warn()
}

func Debug() *zerolog.Event {
	return Logger.Debug()
}

func Error() *zerolog.Event {
	return Logger.Error()
}

func Fatal() *zerolog.Event {
	return Logger.Fatal()
}

func Print(v ...interface{}) {
	Logger.Print(v...)
}

func Printf(format string, v ...interface{}) {
	Logger.Printf(format, v...)
}
