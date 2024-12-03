package notifications

import (
	"bytes"
	"fmt"
	"html/template"
	"net/url"
	"path/filepath"
	"strconv"

	"github.com/john-vh/college_testing/backend/models"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
	if err != nil {
		return "", err
	}

	var res bytes.Buffer
	err = t.Execute(&res, data)
	if err != nil {
		return "", err
	}

	return res.String(), nil
}

type ApplicationUpdatedNotification struct {
	recipient    *models.User
	applicant    *models.User
	application  *models.UserApplication
	postURI      string
	templatePath string
}

func (ns *NotificationsService) NewApplicationUpdatedNotification(recipient *models.User, applicant *models.User, application *models.UserApplication) *ApplicationUpdatedNotification {
	templateName := fmt.Sprintf("Application%v", cases.Title(language.English).String(string(application.Status)))
	postURI, _ := url.JoinPath(ns.frontendURL, "businesses", application.Business.Id.String(), "posts", strconv.Itoa(application.Post.Id))
	return &ApplicationUpdatedNotification{
		recipient:   recipient,
		applicant:   applicant,
		application: application,
		// FIXME: Need to link to the post
		postURI:      postURI,
		templatePath: filepath.Join(ns.templatesPath, templateName) + ".html",
	}
}

func (n *ApplicationUpdatedNotification) To() *models.User {
	return n.applicant
}

func (n *ApplicationUpdatedNotification) Subject() string {
	return fmt.Sprintf("Application %s", n.applicant.Status)
}

func (n *ApplicationUpdatedNotification) HTML() (string, error) {
	type templateData struct {
		RecipientName string
		ApplicantName string
		PostName      string
		PostLink      string
	}

	data := templateData{
		RecipientName: n.recipient.Name,
		ApplicantName: n.applicant.Name,
		PostName:      n.application.Post.Title,
		PostLink:      n.postURI,
	}

	t, err := template.ParseFiles(n.templatePath)
	if err != nil {
		return "", err
	}

	var res bytes.Buffer
	err = t.Execute(&res, data)
	if err != nil {
		return "", err
	}

	return res.String(), nil
}
