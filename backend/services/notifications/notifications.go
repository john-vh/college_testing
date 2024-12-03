package notifications

import "github.com/john-vh/college_testing/backend/models"

type ApplicationReceivedNotification struct {
	users []*models.User
}

func NewApplicationReceivedNotification(users []*models.User) *ApplicationReceivedNotification {
	return &ApplicationReceivedNotification{users: users}
}

func (n *ApplicationReceivedNotification) HTML(u *models.User) string {
	return "<h1>Application Received</h1>"
}

func (n *ApplicationReceivedNotification) Subject() string {
	return "Application Received"
}

func (n *ApplicationReceivedNotification) Targets() []*models.User {
	return n.users
}
