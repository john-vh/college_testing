package notifications

import (
	"context"
	"log/slog"
	"time"

	"github.com/john-vh/college_testing/backend/models"
)

type NotificationsService struct {
	logger        *slog.Logger
	dataStream    chan Notification
	mailClient    *MailClient
	templatesPath string
	frontendURL   string
}

type Notification interface {
	To() *models.User
	Subject() string
	HTML() (string, error)
	ShouldNotify() bool
}

func NewNotificationService(mailClient *MailClient, frontendURL, templatesPath string, logger *slog.Logger) *NotificationsService {
	const notificationBufferSize = 8
	return &NotificationsService{
		dataStream:    make(chan Notification, notificationBufferSize),
		templatesPath: templatesPath,
		frontendURL:   frontendURL,
		logger:        logger,
		mailClient:    mailClient,
	}
}

// Background service interface implementations
func (ns *NotificationsService) Start() {
	go ns.run()
}

func (ns *NotificationsService) Stop() {
	close(ns.dataStream)
}

func (ns *NotificationsService) run() {
	for noti := range ns.dataStream {
		if !noti.ShouldNotify() {
			continue
		}
		user := noti.To()
		body, err := noti.HTML()
		if err != nil {
			ns.logger.Warn("Failed to parse body of notification", "err", err)
			return
		}
		err = ns.mailClient.SendMsg(
			[]string{user.Email},
			&MailInfo{
				ToList:  []string{user.Email},
				Subject: noti.Subject(),
				Body:    body,
			})
		if err != nil {
			ns.logger.Warn("Failed to send mail message", "err", err)
			return
		}
	}
}

func (ns *NotificationsService) EnqueueWithTimeout(ctx context.Context, n Notification) error {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*100)
	defer cancel()
	return ns.Enqueue(ctx, n)
}

func (ns *NotificationsService) Enqueue(ctx context.Context, n Notification) error {
	select {
	case ns.dataStream <- n:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}
