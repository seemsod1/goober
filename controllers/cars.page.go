package controllers

import (
	models "help/models/app_models"
	"help/models/entities"
	"net/http"
)

func CarsPage(w http.ResponseWriter, r *http.Request) {
	var cities []entities.City

	var cars []entities.Car
	result := app.DB.Find(&cities)
	if result.Error != nil {
		http.Error(w, "Failed to fetch cities data", http.StatusInternalServerError)
		return
	}
	if err := app.DB.Preload("Model").Preload("Model.Brand").Preload("Type").Preload("Location").Preload("Location.City").Find(&cars).Error; err != nil {
		http.Error(w, "Failed to fetch locations data", http.StatusInternalServerError)
		return
	}

	var types []entities.CarType
	if err := app.DB.Find(&types).Error; err != nil {
		http.Error(w, "Failed to fetch car types data", http.StatusInternalServerError)
		return
	}

	data := &models.TemplateData{
		Data: map[string]interface{}{
			"Cities": cities,
			"Cars":   cars,
			"Types":  types,
		},
	}

	renderTemplate(w, r, "cars.page.tmpl", data)
}
