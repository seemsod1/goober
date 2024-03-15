package models

import "html/template"

// MailData holds an email message
type MailData struct {
	To      string
	From    string
	Subject string
	Content template.HTML
}
