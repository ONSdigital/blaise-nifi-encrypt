package util

import (
	"os"
	"strings"
	"time"

	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	LogFormat string = "LOG_FORMAT"
	Terminal  string = "Terminal"
	Json      string = "Json"
	Debug     string = "DEBUG"
)

func ConfigureLogging() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// change log format
	if terminal, isFound := os.LookupEnv(LogFormat); isFound {
		switch terminal {
		case Terminal:
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: colorable.NewColorableStdout(), NoColor: false, TimeFormat: time.Stamp})
		case Json:
			// json is the default
		}
	}

	if debug, f := os.LookupEnv(Debug); f {
		switch strings.ToLower(debug) {
		case "true":
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		}
	}
}
