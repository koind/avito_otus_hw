package app

import "github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/domain/entity"

type NotificationSource interface {
	GetNotificationChannel() (<-chan entity.Notification, error)
}

type NotificationTransport interface {
	String() string
	Send(entity.Notification) error
}

type NotificationSender struct {
	source NotificationSource
	logger Logger
}

func NewSender(source NotificationSource, logger Logger) *NotificationSender {
	return &NotificationSender{source, logger}
}

func (s *NotificationSender) Run() {
	s.logger.Info("[notification] Run")

	channel, err := s.source.GetNotificationChannel()
	if err != nil {
		s.logger.Error("[notification] Error get from channel: %s", err)
		return
	}

	for notification := range channel {
		s.logger.Info("[notification] %s", notification)
	}
}
