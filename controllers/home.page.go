package controllers

import (
	"help/helpers/render"
	models "help/models/app_models"
	"net/http"
)

func (m *Repository) HomePage(w http.ResponseWriter, r *http.Request) {
	m.ClearSessionData(r)

	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	m.ClearSessionData(r)

	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again."

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
