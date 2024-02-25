package initializers

import (
	"github.com/justinas/nosurf"
	models "help/models/app_models"
	"net/http"
)

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)

	return td
}
