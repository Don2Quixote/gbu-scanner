package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// NewLogrus returns setted up logrus logger.
func NewLogrus() *logrus.Logger {
	return &logrus.Logger{
		Out:   os.Stdout,
		Level: logrus.DebugLevel,
		Formatter: &logrus.TextFormatter{
			FullTimestamp: true,
			ForceColors:   true,
		},
	}
}
