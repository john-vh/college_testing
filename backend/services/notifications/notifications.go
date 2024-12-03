package notifications

import (
	"bytes"
	"html/template"
	"net/url"
	"path/filepath"
	"strconv"

	"github.com/john-vh/college_testing/backend/models"
)

type ApplicationReceivedNotification struct {
	recipient    *models.User
	applicant    *models.User
	post         *models.Post
	postURI      string
	templatePath string
}

func (ns *NotificationsService) NewApplicationReceivedNotification(recipient *models.User, applicant *models.User, post *models.Post) *ApplicationReceivedNotification {
	const templateName = "ApplicationReceived"
	postURI, _ := url.JoinPath(ns.frontendURL, "businesses", post.BusinessId.String(), "posts", strconv.Itoa(post.Id))
	return &ApplicationReceivedNotification{
		recipient: recipient,
		applicant: applicant,
		post:      post,
		// FIXME: Need to link to the post
		postURI:      postURI,
		templatePath: filepath.Join(ns.templatesPath, templateName) + ".html",
	}
}

func (n *ApplicationReceivedNotification) To() *models.User {
	return n.recipient
}

func (n *ApplicationReceivedNotification) Subject() string {
	return "Application Received"
}

func (n *ApplicationReceivedNotification) HTML() (string, error) {
	type templateData struct {
		Name          string
		ApplicantName string
		PostName      string
		PostLink      string
	}

	data := templateData{
		Name:          n.recipient.Name,
		ApplicantName: n.applicant.Name,
		PostName:      n.post.Title,
		PostLink:      n.postURI,
	}

	t, err := template.ParseFiles(n.templatePath)

	var res bytes.Buffer
	if err != nil {
		return "", err
	}
	err = t.Execute(&res, data)
	if err != nil {
		return "", err
	}

	return res.String(), nil
}
