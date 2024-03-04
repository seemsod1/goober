package controllers

import (
	"encoding/json"
	"help/helpers"
	"help/helpers/render"
	models "help/models/app_models"
	"help/models/entities"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (m *Repository) AllCars(w http.ResponseWriter, r *http.Request) {
	userId, _ := m.App.Session.Get(r.Context(), "user_id").(int)

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	perPage := 9

	var location entities.RentLocation
	res := m.App.DB.Table("rent_locations").Preload("City").Where("user_id = ?", userId).Take(&location)
	if res.Error != nil {
		m.App.Session.Put(r.Context(), "error", "Can't get locations")
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	var cars []entities.Car

	var totalItems int64
	res = m.App.DB.Table("cars").Where("location_id = ?", location.ID).Count(&totalItems)

	offset := (page - 1) * perPage

	res = m.App.DB.Table("cars").Preload("Model").Preload("Model.Brand").Preload("Type").Where("location_id = ?", location.ID).Offset(offset).Limit(perPage).Find(&cars)
	if res.Error != nil {
		m.App.Session.Put(r.Context(), "error", "Can't get cars")
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	startDate, err := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse start date!")
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	var availableCars []entities.Car
	res = m.App.DB.
		Debug().
		Preload("Location").
		Preload("Location.City").
		Preload("Model").
		Preload("Model.Brand").
		Preload("Type").
		Table("cars").
		Select("cars.*").
		Joins("JOIN rent_locations ON cars.location_id = rent_locations.id").
		Joins("JOIN cities ON rent_locations.city_id = cities.id").
		Where("cities.name = ?", location.City.Name).
		Not("cars.id IN (?)", m.App.DB.Table("car_histories").
			Select("DISTINCT car_id").
			Joins("JOIN rent_infos ON car_histories.rent_info_id = rent_infos.id").
			Where("rent_infos.status_id IN (?)", []int{4, 1}).
			Where(
				"(? BETWEEN DATE(rent_infos.start_date) AND DATE(rent_infos.end_date)) OR "+
					"(? BETWEEN DATE(rent_infos.start_date) AND DATE(rent_infos.end_date)) OR "+
					"(DATE(rent_infos.start_date) <= ? AND DATE(rent_infos.end_date) >= ?)", startDate, startDate, startDate, startDate),
		).
		Find(&availableCars)
	if res.Error != nil {
		for _, car := range cars {
			car.Available = false
		}
	}
	for id, car := range cars {
		if findCar(availableCars, car.ID) {
			cars[id].Available = true
		} else {
			cars[id].Available = false
		}
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(perPage)))

	data := make(map[string]interface{})
	data["cars"] = cars
	data["pagination"] = models.PaginationData{
		CurrentPage: page,
		TotalPages:  totalPages,
		PrevPage:    page - 1,
		NextPage:    page + 1,
		HasPrev:     page > 1,
		HasNext:     page < totalPages,
		Pages:       helpers.GeneratePageNumbers(page, totalPages),
	}

	render.RenderTemplate(w, r, "head-all-cars.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func findCar(cars []entities.Car, id int) bool {
	for _, car := range cars {
		if car.ID == id {
			return true
		}
	}
	return false
}

func (m *Repository) CarHistory(w http.ResponseWriter, r *http.Request) {

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	perPage := 5

	exploded := strings.Split(r.RequestURI, "/")
	carq := strings.Split(exploded[3], "?")
	carId, err := strconv.Atoi(carq[0])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "missing url parameter")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	offset := (page - 1) * perPage

	var carHistories []entities.CarHistory

	var totalItems int64
	m.App.DB.Table("car_histories").Where("car_id = ?", carId).Count(&totalItems)

	res := m.App.DB.Preload("RentInfo").Preload("RentInfo.Status").Preload("RentInfo.From").Preload("RentInfo.From.City").Preload("RentInfo.Return").Preload("RentInfo.Return.City").Preload("Car").Preload("Car.Model").Preload("Car.Model.Brand").Table("car_histories").Where("car_id = ?", carId).Offset(offset).Limit(perPage).Order("created_at desc").Find(&carHistories)
	if res.Error != nil {
		m.App.Session.Put(r.Context(), "error", "Can't get car histories")
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(perPage)))

	responseData := map[string]interface{}{
		"carHistories": carHistories,
		"pagination": models.PaginationData{
			CurrentPage: page,
			TotalPages:  totalPages,
			PrevPage:    page - 1,
			NextPage:    page + 1,
			HasPrev:     page > 1,
			HasNext:     page < totalPages,
			Pages:       helpers.GeneratePageNumbers(page, totalPages),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}
