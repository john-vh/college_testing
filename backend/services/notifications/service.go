package notifications

import (
	"context"
	"log/slog"
	"time"

	"github.com/john-vh/college_testing/backend/models"
	"github.com/john-vh/college_testing/backend/util"
)

type NotificationsService struct {
	logger     *slog.Logger
	dataStream chan Notification
	mailClient *MailClient
}

type Notification interface {
	Targets() []*models.User
	Subject() string
	HTML(*models.User) string
}

func NewNotificationService(mailClient *MailClient, logger *slog.Logger) *NotificationsService {
	const notificationBufferSize = 8
	return &NotificationsService{
		logger:     logger,
		dataStream: make(chan Notification, notificationBufferSize),
		mailClient: mailClient,
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
		filtered := util.Filter(noti.Targets(), func(u *models.User) bool { return ns.shouldNotify(u, noti) })

		for _, user := range filtered {
			go func(user *models.User, noti Notification) {
				ns.mailClient.SendMsg(
					[]string{user.Email},
					&MailInfo{
						ToList:  []string{user.Email},
						Subject: noti.Subject(),
						Body:    noti.HTML(user),
					})
			}(user, noti)
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

func (ns *NotificationsService) shouldNotify(user *models.User, n Notification) bool {
	switch n.(type) {
	case *ApplicationReceivedNotification:
		return true
	}

	ns.logger.Warn("Uncaught notification type", "notification", n)

	return false

}
