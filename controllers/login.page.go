package controllers

import (
	models "help/models/app_models"
	"net/http"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, "login.page.tmpl", &models.TemplateData{})
}
