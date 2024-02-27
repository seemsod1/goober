package controllers

import (
	"help/helpers/render"
	models "help/models/app_models"
	"net/http"
)

func (m *Repository) LoginPage(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "login.page.tmpl", &models.TemplateData{})
}
