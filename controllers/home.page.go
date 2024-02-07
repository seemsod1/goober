package controllers

import (
	models "help/models/app_models"
	"net/http"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	app.Session.Put(r.Context(), "remote_ip", remoteIP)

	renderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

func About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"

	remoteIp := app.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIp

	renderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
