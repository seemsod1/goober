package controllers

import (
	"help/helpers/render"
	models "help/models/app_models"
	"net/http"
)

func (m *Repository) NotFoundPage(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "404.page.tmpl", &models.TemplateData{})
}
