package controllers

import (
	"help/helpers/render"
	models "help/models/app_models"
	"help/models/entities"
	"net/http"
)

func (m *Repository) CarsPage(w http.ResponseWriter, r *http.Request) {
	var cities []entities.City

	var cars []entities.Car
	result := m.App.DB.Find(&cities)
	if result.Error != nil {
		http.Error(w, "Failed to fetch cities data", http.StatusInternalServerError)
		return
	}
	if err := m.App.DB.Preload("Model").Preload("Model.Brand").Preload("Type").Preload("Location").Preload("Location.City").Find(&cars).Error; err != nil {
		http.Error(w, "Failed to fetch locations data", http.StatusInternalServerError)
		return
	}

	var types []entities.CarType
	if err := m.App.DB.Find(&types).Error; err != nil {
		http.Error(w, "Failed to fetch car types data", http.StatusInternalServerError)
		return
	}

	// Fetch car assignments
	var assignments []entities.CarAssignment
	if err := m.App.DB.Preload("Purpose").Preload("Car").Find(&assignments).Error; err != nil {
		http.Error(w, "Failed to fetch car assignments data", http.StatusInternalServerError)
		return
	}

	// Group assignments by car ID
	assignmentsByCarID := make(map[int][]entities.CarAssignment)
	for _, assignment := range assignments {
		carID := assignment.Car.ID
		assignmentsByCarID[carID] = append(assignmentsByCarID[carID], assignment)
	}

	// Associate assignments with their respective cars
	for i, _ := range cars {
		cars[i].Assignments = assignmentsByCarID[cars[i].ID]
	}

	// Pass data to the HTML template
	data := &models.TemplateData{
		Data: map[string]interface{}{
			"Cities": cities,
			"Cars":   cars,
			"Types":  types,
			// No need to pass assignments separately as they are already associated with cars
		},
	}

	render.RenderTemplate(w, r, "cars.page.tmpl", data)
}
