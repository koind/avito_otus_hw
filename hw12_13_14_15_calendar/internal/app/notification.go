package app

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	EventID uuid.UUID
	Title   string
	Dt      time.Time
	UserID  string
}

func (n Notification) String() string {
	builder := strings.Builder{}
	builder.WriteString("Notification: ")
	builder.WriteString(n.Title)
	builder.WriteString(" at ")
	builder.WriteString(n.Dt.Format(time.RFC3339))

	return builder.String()
}
