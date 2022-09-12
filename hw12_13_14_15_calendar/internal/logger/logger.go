package logger

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	logger *logrus.Logger
}

func New(loggerConfig config.LoggerConf) (*Logger, error) {
	logger := logrus.New()

	loggerOutput, err := parseLogFile(loggerConfig.Filename)
	if err != nil {
		return nil, fmt.Errorf("invalid log file name: %w", err)
	}
	logger.SetOutput(loggerOutput)

	loggerLevel, err := logrus.ParseLevel(string(loggerConfig.Level))
	if err != nil {
		return nil, err
	}
	logger.SetLevel(loggerLevel)

	logger.SetFormatter(&logrus.JSONFormatter{})

	return &Logger{
		logger,
	}, nil
}

func (l *Logger) Debug(message string, params ...interface{}) {
	l.logger.Debugf(message, params...)
}

func (l *Logger) Info(message string, params ...interface{}) {
	l.logger.Infof(message, params...)
}

func (l *Logger) Warn(message string, params ...interface{}) {
	l.logger.Warnf(message, params...)
}

func (l *Logger) Error(message string, params ...interface{}) {
	l.logger.Errorf(message, params...)
}

func (l *Logger) LogHTTPRequest(r *http.Request, code, length int) {
	l.logger.Infof(
		"%s [%s] %s %s %s %d %d %q",
		r.RemoteAddr,
		time.Now().Format("02/Jan/2006:15:04:05 MST"),
		r.Method,
		r.RequestURI,
		r.Proto,
		code,
		length,
		r.UserAgent(),
	)
}

func parseLogFile(filename string) (io.Writer, error) {
	switch filename {
	case "stderr":
		return os.Stderr, nil
	case "stdout":
		return os.Stdout, nil
	default:
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
		if err != nil {
			return nil, err
		}
		return file, nil
	}
}
