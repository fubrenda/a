package logzer

import (
	"io"
	"log"
	stdliblog "log"
	"os"
	"runtime"

	"github.com/rs/zerolog"
)

var logger zerolog.Logger

// MustNewLogzer setups a new zerologger
func MustNewLogzer(name string, verbose bool, out io.Writer) zerolog.Logger {
	if verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Setup a logger.
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("failed to get hostname")
	}

	if out == nil {
		out = os.Stdout
	}

	logger = zerolog.New(out).With().
		Timestamp().
		Str("name", name).
		Str("host", hostname).
		Str("go_version", runtime.Version()).
		Logger()

	// Tell stdlib tools that use the default logger to use our logger.
	stdliblog.SetFlags(0)
	stdliblog.SetOutput(logger)

	return logger
}
