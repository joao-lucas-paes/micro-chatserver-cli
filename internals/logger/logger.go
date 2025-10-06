package logger

import (
	"io"
	"os"
	"sync"
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
type Logger struct {
	logger log.Logger
	mu sync.Mutex
}

func New(logfile string) (Logger, error) {
	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return Logger{}, err
	}

	std := io.MultiWriter(os.Stdout, file)

	return Logger{
		logger: *log.NewWithOptions(std, log.Options{
			ReportCaller:    false,
			ReportTimestamp: true,
		}),
		mu: sync.Mutex{},
	}, nil
}



func (l *Logger) Println(message string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Print(message, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Infof(format, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Warnf(format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Errorf(format, args...)
}
