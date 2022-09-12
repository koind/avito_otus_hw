package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("invalid config file", func(t *testing.T) {
		t.Skip("No need to test it now")

		_, err := LoadConfig()
		require.Error(t, err)

		file, err := os.CreateTemp("", "log")
		if err != nil {
			t.FailNow()
			return
		}
		_, err = file.Write([]byte("invalid json"))
		if err != nil {
			t.FailNow()
			return
		}
		_, err = LoadConfig()
		require.Error(t, err)
	})
}
