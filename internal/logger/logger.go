package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

var Log *zerolog.Logger

func init() {
	if Log != nil {
		return
	}

	l := zerolog.New(os.Stdout).With().Logger()
	Log = &l
}

func Out(w io.Writer) {
	l := zerolog.New(w).With().Logger()
	Log = &l
}
