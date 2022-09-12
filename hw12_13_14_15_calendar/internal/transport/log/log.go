package transport

import "github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/app"

type LogNotificationTransport struct {
	logger app.Logger
}

func NewLogNotificationTransport(logger app.Logger) *LogNotificationTransport {
	return &LogNotificationTransport{logger: logger}
}

func (t *LogNotificationTransport) String() string {
	return "LogNotificationTransport"
}

func (t *LogNotificationTransport) Send(n app.Notification) error {
	t.logger.Info("[notification][transport][log] %s", n)
	return nil
}
