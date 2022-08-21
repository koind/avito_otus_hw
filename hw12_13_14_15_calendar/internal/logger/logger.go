package logger

import (
	"fmt"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
)

type Logger struct {
	logger *logrus.Logger
}

func New(cfg config.LoggerConf) (*Logger, error) {
	log := logrus.New()

	output, err := openLog(cfg.Filename)
	if err != nil {
		return nil, fmt.Errorf("error open log file: %w", err)
	}

	log.SetOutput(output)

	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	log.SetLevel(level)
	log.SetFormatter(&logrus.JSONFormatter{})

	return &Logger{log}, nil
}

func (l *Logger) Info(message string, params ...interface{}) {
	l.logger.Infof(message, params...)
}

func (l *Logger) Error(message string, params ...interface{}) {
	l.logger.Errorf(message, params...)
}

func (l *Logger) LogRequest(r *http.Request, code, length int) {
	l.logger.Infof(
		"%s %s %s %s %d %d %q",
		r.RemoteAddr,
		r.Method,
		r.RequestURI,
		r.Proto,
		code,
		length,
		r.UserAgent(),
	)
}

func openLog(file string) (io.Writer, error) {
	switch file {
	case "stderr":
		fmt.Println("stderr")
		return os.Stderr, nil
	case "stdout":
		fmt.Println("stdout")
		return os.Stdout, nil
	default:
		file, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
		if err != nil {
			return nil, err
		}
		return file, nil
	}
}
