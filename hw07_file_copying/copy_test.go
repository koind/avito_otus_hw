package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("Test /dev/urandom", func(t *testing.T) {
		err := Copy("/dev/urandom", "/tmp/", 0, 0)
		require.Equal(t, ErrUnsupportedFile, err)
	})

	t.Run("Test exceeded offset", func(t *testing.T) {
		tempFile, err := ioutil.TempFile(".", "out_")
		if err != nil {
			t.FailNow()
		}
		defer func() {
			tempFile.Close()
			os.Remove(tempFile.Name())
		}()

		err = Copy("./testdata/input.txt", tempFile.Name(), 10000, 0)
		require.Equal(t, ErrOffsetExceedsFileSize, err)
	})

	t.Run("Non-existent file opening", func(t *testing.T) {
		err := Copy("./lsdkjflajflksdjfls;fk", "doesn't_matter.jpg", 0, 0)
		require.Error(t, err)
	})
}
