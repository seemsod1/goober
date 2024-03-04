package controllers

import (
	"encoding/json"
	"help/helpers/forms"
	"help/helpers/render"
	models "help/models/app_models"
	"help/models/entities"
	"net/http"
	"strconv"
	"time"
)

func (m *Repository) AddCar(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "head-add-car.page.tmpl", &models.TemplateData{})
}

func (m *Repository) GetBrands(w http.ResponseWriter, r *http.Request) {
	var brands []entities.CarBrand

	result := m.App.DB.
		Distinct().
		Table("car_brands").
		Joins("JOIN car_models ON car_brands.id = car_models.brand_id").
		Find(&brands)

	if result.Error != nil {
		http.Error(w, "Failed to fetch brands data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(brands)
}

func (m *Repository) GetModels(w http.ResponseWriter, r *http.Request) {
	brandId := r.URL.Query().Get("brandId")
	var mod []entities.CarModel

	result := m.App.DB.Table("car_models").Where("brand_id = ?", brandId).Find(&mod)
	if result.Error != nil {
		http.Error(w, "Failed to fetch models data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mod)
}

func (m *Repository) GetTypes(w http.ResponseWriter, r *http.Request) {
	var types []entities.CarType

	result := m.App.DB.Find(&types)
	if result.Error != nil {
		http.Error(w, "Failed to fetch types data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(types)
}

func (m *Repository) AddCarPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Can't parse form", http.StatusBadRequest)
		return
	}
	model := r.Form.Get("SelectModel")
	ctype := r.Form.Get("SelectType")
	bags := r.Form.Get("inputBags")
	passengers := r.Form.Get("inputPassengers")
	year := r.Form.Get("year")
	price := r.Form.Get("price")
	plate := r.Form.Get("plate")
	color := r.Form.Get("color")

	form := forms.New(r.PostForm)
	form.Required("SelectBrand", "SelectModel", "SelectType", "inputBags", "inputPassengers", "year", "price", "plate", "color")
	form.IsNumber("inputBags")
	form.IsNumber("inputPassengers")
	form.IsNumber("year")
	form.IsNumber("price")
	form.IsPlate("plate")

	if !form.Valid() {
		render.RenderTemplate(w, r, "head-add-car.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}
	modelId, err := strconv.Atoi(model)
	if err != nil || modelId < 1 {
		m.App.Session.Put(r.Context(), "error", "Invalid model id")
		return
	}
	typeId, err := strconv.Atoi(ctype)
	if err != nil || typeId < 1 {
		m.App.Session.Put(r.Context(), "error", "Invalid type id")
		return
	}
	bagsInt, err := strconv.Atoi(bags)
	if err != nil || bagsInt < 1 || bagsInt > 6 {
		m.App.Session.Put(r.Context(), "error", "Invalid bags")
		return
	}
	passengersInt, err := strconv.Atoi(passengers)
	if err != nil || passengersInt < 2 || passengersInt > 6 {
		m.App.Session.Put(r.Context(), "error", "Invalid passengers")
		return
	}

	yearInt, err := strconv.Atoi(year)
	if err != nil || yearInt < 1950 || yearInt > time.Now().Year() {
		m.App.Session.Put(r.Context(), "error", "Invalid year")
		return
	}
	priceInt, err := strconv.Atoi(price)
	if err != nil || priceInt < 1 || priceInt > 10000 {
		m.App.Session.Put(r.Context(), "error", "Invalid price")
		return
	}

	userId := m.App.Session.GetInt(r.Context(), "user_id")

	var location entities.RentLocation
	result := m.App.DB.Table("rent_locations").Where("user_id = ?", userId).Take(&location)
	if result.Error != nil {
		m.App.Session.Put(r.Context(), "error", "Location not found")
		return
	}

	var car entities.Car
	car = entities.Car{
		ModelId:    modelId,
		TypeId:     typeId,
		Bags:       bagsInt,
		Passengers: passengersInt,
		Year:       yearInt,
		Price:      float64(priceInt),
		Plate:      plate,
		Color:      color,
		LocationId: location.ID,
	}

	res := m.App.DB.Create(&car)
	if res.Error != nil {
		m.App.Session.Put(r.Context(), "error", "Can't create car")
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Successfully Added!")
	http.Redirect(w, r, "/head/add-car", http.StatusSeeOther)
}
