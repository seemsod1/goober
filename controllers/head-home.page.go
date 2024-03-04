package controllers

import (
	"help/helpers/render"
	models "help/models/app_models"
	"net/http"
)

func (m *Repository) HeadPage(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "head.page.tmpl", &models.TemplateData{})
}
