// Trying to avoid creating a `utils` or `helpers` package

package logging

import (
	"strings"

	"github.com/sirupsen/logrus"
)

// Create a new logger with the given log level
// If localDev is true, the logger will output human readable logs
// Otherwise, it will output structured logs
func NewLogger(level string, localDev bool) *logrus.Entry {
	var format string = "2006-01-02 15:04:05"
	var formatter logrus.Formatter = &logrus.JSONFormatter{TimestampFormat: format}
	var logLevel logrus.Level

	switch strings.ToUpper(level) {
	case "DEBUG":
		logLevel = logrus.DebugLevel
	case "WARN":
		logLevel = logrus.WarnLevel
	case "ERROR":
		logLevel = logrus.ErrorLevel
	case "FATAL":
		logLevel = logrus.FatalLevel
	default:
		logLevel = logrus.InfoLevel
	}

	// If localDev is true, we want to output human readable logs
	if localDev {
		formatter = &logrus.TextFormatter{
			ForceColors:     true,
			FullTimestamp:   true,
			TimestampFormat: format,
		}
	}

	log := logrus.New()
	log.SetFormatter(formatter)
	log.SetLevel(logLevel)
	return log.WithField("app", "tomedome")
}
