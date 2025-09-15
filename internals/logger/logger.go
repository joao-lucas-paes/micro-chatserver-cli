package logger

import (
	"io"
	"os"
	"github.com/charmbracelet/log"
)

// Logger defines a generic logging interface with methods for printing messages
// and formatted logs at different severity levels.
//
// The interface provides the following methods:
//   - Println: Logs a message with a newline.
//   - Infof: Logs an informational message with formatting.
//   - Warnf: Logs a warning message with formatting.
//   - Errorf: Logs an error message with formatting.
type Logger interface {
	Println(message string)
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

func New(logfile string) (Logger, error) {
	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	std := io.MultiWriter(os.Stdout, file)

	return &defaultLogger{
		logger: *log.NewWithOptions(std, log.Options{
			ReportCaller:    false,
			ReportTimestamp: true,
		}),
	}, nil
}

type defaultLogger struct {
	logger log.Logger
}

func (l *defaultLogger) Println(message string) {
	l.logger.Print(message)
}

func (l *defaultLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *defaultLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *defaultLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}
