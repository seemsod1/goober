package controllers

import (
	"github.com/seemsod1/goober/helpers/render"
	models "github.com/seemsod1/goober/models/app_models"
	"net/http"
)

func (m *Repository) NotFoundPage(w http.ResponseWriter, r *http.Request) {
	m.ClearSessionData(r)

	render.RenderTemplate(w, r, "404.page.tmpl", &models.TemplateData{})
}
