package business

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/john-vh/college_testing/backend/models"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type businessRequestedAdminNotification struct {
	recipient    *models.User
	applicant    *models.User
	business     *models.Business
	businessURI  string
	templatePath string
}

func (h *BusinessHandler) NewBusinessRequestedAdminNotification(recipeint, applicant *models.User, b *models.Business) *businessRequestedAdminNotification {
	const templateName = "BusinessRequestedAdmin"
	// FIXME: Ignoring error
	URI, _ := b.URI(h.frontendURL)
	return &businessRequestedAdminNotification{
		recipient:    recipeint,
		applicant:    applicant,
		business:     b,
		businessURI:  URI,
		templatePath: filepath.Join(h.notificationsTemplatesDir, templateName) + ".html",
	}
}

func (n *businessRequestedAdminNotification) ShouldNotify() bool { return true }
func (n *businessRequestedAdminNotification) To() *models.User   { return n.recipient }
func (n *businessRequestedAdminNotification) Subject() string    { return "New Business Requested" }
func (n *businessRequestedAdminNotification) HTML() (string, error) {
	type templateData struct {
		RecipientName string
		ApplicantName string
		BusinessName  string
		BusinessURI   string
	}

	data := templateData{
		RecipientName: n.recipient.Name,
		ApplicantName: n.applicant.Name,
		BusinessName:  n.business.Name,
		BusinessURI:   n.businessURI,
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

type ApplicationReceivedNotification struct {
	recipient    *models.User
	applicant    *models.User
	post         *models.Post
	postURI      string
	templatePath string
}

func (h *BusinessHandler) NewApplicationReceivedNotification(recipient *models.User, applicant *models.User, post *models.Post) *ApplicationReceivedNotification {
	const templateName = "ApplicationReceived"
	postURI := ""
	return &ApplicationReceivedNotification{
		recipient: recipient,
		applicant: applicant,
		post:      post,
		// FIXME: Need to link to the post
		postURI:      postURI,
		templatePath: filepath.Join(h.notificationsTemplatesDir, templateName) + ".html",
	}
}
func (n *ApplicationReceivedNotification) ShouldNotify() bool {
	return n.To().NotifyApplicationReceived
}
func (n *ApplicationReceivedNotification) To() *models.User { return n.recipient }
func (n *ApplicationReceivedNotification) Subject() string  { return "Application Received" }
func (n *ApplicationReceivedNotification) HTML() (string, error) {
	type templateData struct {
		RecipientName string
		ApplicantName string
		PostName      string
		PostLink      string
	}

	data := templateData{
		RecipientName: n.recipient.Name,
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

type ApplicationWithdrawnNotification struct {
	recipient      *models.User
	applicant      *models.User
	application    *models.UserApplication
	applicationURI string
	templatePath   string
}

func (h *BusinessHandler) NewApplicationWithdrawnNotification(recipient *models.User, applicant *models.User, application *models.UserApplication) *ApplicationWithdrawnNotification {
	templateName := "ApplicationWithdrawn"
	// FIXME: Ignoring error
	URI, _ := application.URI(h.frontendURL)
	return &ApplicationWithdrawnNotification{
		recipient:      recipient,
		applicant:      applicant,
		application:    application,
		applicationURI: URI,
		templatePath:   filepath.Join(h.notificationsTemplatesDir, templateName) + ".html",
	}
}

func (n *ApplicationWithdrawnNotification) ShouldNotify() bool {
	return n.To().NotifyApplicationWithdrawn
}

func (n *ApplicationWithdrawnNotification) To() *models.User {
	return n.recipient
}

func (n *ApplicationWithdrawnNotification) Subject() string {
	return fmt.Sprintf("Application Withdrawn")
}

func (n *ApplicationWithdrawnNotification) HTML() (string, error) {
	type templateData struct {
		RecipientName   string
		ApplicantName   string
		BusinessName    string
		PostName        string
		ApplicationLink string
	}

	data := templateData{
		RecipientName:   n.recipient.Name,
		ApplicantName:   n.applicant.Name,
		BusinessName:    n.application.Business.Name,
		PostName:        n.application.Post.Title,
		ApplicationLink: n.applicationURI,
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
	applicant      *models.User
	application    *models.UserApplication
	applicationURI string
	templatePath   string
}

func (h *BusinessHandler) NewApplicationUpdatedNotification(applicant *models.User, application *models.UserApplication) *ApplicationUpdatedNotification {
	templateName := fmt.Sprintf("Application%v", cases.Title(language.English).String(string(application.Status)))
	// FIXME: Ignoring error
	applicationURI, _ := application.URI(h.frontendURL)
	return &ApplicationUpdatedNotification{
		applicant:      applicant,
		application:    application,
		applicationURI: applicationURI,
		templatePath:   filepath.Join(h.notificationsTemplatesDir, templateName) + ".html",
	}
}

func (n *ApplicationUpdatedNotification) ShouldNotify() bool {
	return n.To().NotifyApplicationUpdated
}

func (n *ApplicationUpdatedNotification) To() *models.User {
	return n.applicant
}

func (n *ApplicationUpdatedNotification) Subject() string {
	return "Application Status Update"
}

func (n *ApplicationUpdatedNotification) HTML() (string, error) {
	type templateData struct {
		ApplicantName   string
		BusinessName    string
		PostName        string
		ApplicationLink string
	}

	data := templateData{
		ApplicantName:   n.applicant.Name,
		BusinessName:    n.application.Business.Name,
		PostName:        n.application.Post.Title,
		ApplicationLink: n.applicationURI,
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
