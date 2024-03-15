package models

import (
	"github.com/alexedwards/scs/v2"
	"github.com/go-co-op/gocron/v2"
	"gorm.io/gorm"
	"html/template"
	"log"
)

// AppConfig holds the application config
type AppConfig struct {
	UseCache      bool
	DB            *gorm.DB
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	Scheduler     gocron.Scheduler
	MailChan      chan MailData
}
