package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/justinas/nosurf"
	models "help/models/app_models"
	"help/models/entities"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

var functions = template.FuncMap{
	"shuffle": func(data string) template.JS {
		return template.JS(data)
	},
	"jsonify": func(data interface{}) template.JS {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return template.JS("")
		}
		return template.JS(jsonData)
	},
	"toLower": func(s string) string {
		return strings.ToLower(s)
	},
}
var app *models.AppConfig

// NewRenderer sets the config for the template package
func NewRenderer(a *models.AppConfig) {
	app = a
}

// AddDefaultData adds data for all templates
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	if app.Session.Exists(r.Context(), "user_role") {
		td.Role = app.Session.GetInt(r.Context(), "user_role")
	}
	if app.Session.GetInt(r.Context(), "user_role") == 2 {
		var location entities.RentLocation
		app.DB.Table("rent_locations").Preload("City").Where("user_id = ?", app.Session.GetInt(r.Context(), "user_id")).Take(&location)
		td.Location = location.City.Name + ", " + location.FullAddress
	}

	td.CSRFToken = nosurf.Token(r)
	return td
}

// RenderTemplate renders a template
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template

	if app.UseCache {
		// get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser", err)
	}

}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./views/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./views/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./views/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}
	return myCache, nil
}
