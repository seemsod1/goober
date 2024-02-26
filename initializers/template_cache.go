package initializers

import (
	"encoding/json"
	"html/template"
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
