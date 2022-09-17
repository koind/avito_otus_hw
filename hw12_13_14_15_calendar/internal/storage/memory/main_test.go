package memorystorage

import (
	"os"
	"testing"

	calendarconfig "github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/configs"
	_ "github.com/lib/pq"
)

var memoryStorage *MemoryStorage

func TestMain(m *testing.M) {
	memoryStorage = New(calendarconfig.Config{})

	code := m.Run()

	os.Exit(code)
}
