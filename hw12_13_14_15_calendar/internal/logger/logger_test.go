package logger

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	t.Run("debug level", func(t *testing.T) {
		file, err := os.CreateTemp("", "log")
		if err != nil {
			t.FailNow()
			return
		}

		defer os.Remove(file.Name())
		defer file.Close()

		l, _ := New(config.LoggerConf{
			Level:    config.Debug,
			Filename: file.Name(),
		})
		l.Debug("DEBUG %s", "debug message")
		l.Info("INFO %s", "info message")
		l.Warn("WARNING %s", "warning message")
		l.Error("ERROR %s", "error message")

		logContent, _ := os.ReadFile(file.Name())
		fmt.Println(string(logContent))

		require.True(t, strings.Contains(string(logContent), "DEBUG debug message"))
		require.True(t, strings.Contains(string(logContent), "INFO info message"))
		require.True(t, strings.Contains(string(logContent), "WARNING warning message"))
		require.True(t, strings.Contains(string(logContent), "ERROR error message"))
	})

	t.Run("info level", func(t *testing.T) {
		file, err := os.CreateTemp("/tmp", "log")
		if err != nil {
			t.FailNow()
			return
		}

		defer os.Remove(file.Name())
		defer file.Close()

		l, _ := New(config.LoggerConf{
			Level:    config.Info,
			Filename: file.Name(),
		})
		l.Debug("DEBUG %s", "debug message")
		l.Info("INFO %s", "info message")
		l.Warn("WARNING %s", "warning message")
		l.Error("ERROR %s", "error message")

		logContent, _ := os.ReadFile(file.Name())
		fmt.Println(string(logContent))

		require.False(t, strings.Contains(string(logContent), "DEBUG debug message"))
		require.True(t, strings.Contains(string(logContent), "INFO info message"))
		require.True(t, strings.Contains(string(logContent), "WARNING warning message"))
		require.True(t, strings.Contains(string(logContent), "ERROR error message"))
	})

	t.Run("warning level", func(t *testing.T) {
		file, err := os.CreateTemp("/tmp", "log")
		if err != nil {
			t.FailNow()
			return
		}

		defer os.Remove(file.Name())
		defer file.Close()

		l, _ := New(config.LoggerConf{
			Level:    config.Warn,
			Filename: file.Name(),
		})
		l.Debug("DEBUG %s", "debug message")
		l.Info("INFO %s", "info message")
		l.Warn("WARNING %s", "warning message")
		l.Error("ERROR %s", "error message")

		logContent, _ := os.ReadFile(file.Name())
		fmt.Println(string(logContent))

		require.False(t, strings.Contains(string(logContent), "DEBUG debug message"))
		require.False(t, strings.Contains(string(logContent), "INFO info message"))
		require.True(t, strings.Contains(string(logContent), "WARNING warning message"))
		require.True(t, strings.Contains(string(logContent), "ERROR error message"))
	})

	t.Run("error level", func(t *testing.T) {
		file, err := os.CreateTemp("/tmp", "log")
		if err != nil {
			t.FailNow()
			return
		}

		defer os.Remove(file.Name())
		defer file.Close()

		l, _ := New(config.LoggerConf{
			Level:    config.Error,
			Filename: file.Name(),
		})
		l.Debug("DEBUG %s", "debug message")
		l.Info("INFO %s", "info message")
		l.Warn("WARNING %s", "warning message")
		l.Error("ERROR %s", "error message")

		logContent, _ := os.ReadFile(file.Name())
		fmt.Println(string(logContent))

		require.False(t, strings.Contains(string(logContent), "DEBUG debug message"))
		require.False(t, strings.Contains(string(logContent), "INFO info message"))
		require.False(t, strings.Contains(string(logContent), "WARNING warning message"))
		require.True(t, strings.Contains(string(logContent), "ERROR error message"))
	})
}
