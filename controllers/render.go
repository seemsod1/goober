package controllers

import (
	"bytes"
	"help/initializers"
	models "help/models/app_models"
	"html/template"
	"log"
	"net/http"
)

var functions = template.FuncMap{}

var app *models.AppConfig

func SetAppForTemplate(a *models.AppConfig) {
	app = a
}

// renderTemplate is a function that renders a template
func renderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	tc := app.TemplateCache

	if app.UseCache {
		// get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = initializers.CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	_ = t.Execute(buf, td)

	td = initializers.AddDefaultData(td)

	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}
