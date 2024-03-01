package controllers

import (
	"help/helpers/render"
	models "help/models/app_models"
	"help/models/entities"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (m *Repository) CarsPage(w http.ResponseWriter, r *http.Request) {
	var cities []entities.City

	result := m.App.DB.Find(&cities)
	if result.Error != nil {
		http.Error(w, "Failed to fetch cities data", http.StatusInternalServerError)
		return
	}
	data := &models.TemplateData{
		Data: map[string]interface{}{
			"Cities": cities,
		},
	}

	render.RenderTemplate(w, r, "cars.page.tmpl", data)
}

func (m *Repository) CarsPagePost(w http.ResponseWriter, r *http.Request) {
	rent, ok := m.App.Session.Get(r.Context(), "rent").(entities.RentInfo)
	if ok {
		m.App.Session.Remove(r.Context(), "rent")
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Can't parse form", http.StatusBadRequest)
		return
	}

	city := r.Form.Get("citySearchInput")
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse start date!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse end date!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var cars []entities.Car

	//result := m.App.DB.Raw("SELECT    c.ID,c.model_id,c.type_id,c.location_id,    c.Plate,    c.Price,    c.Color,    c.Bags,    c.Passengers,    c.Year FROM cars c        LEFT JOIN car_histories ch ON c.ID = ch.Car_Id        LEFT JOIN rent_infos r ON ch.rent_info_id = r.ID         LEFT JOIN cities l ON c.location_id = l.ID WHERE l.name = ?  AND (r.ID IS NULL OR (r.Start_Date > ? OR r.End_Date < ?))ORDER BY c.ID", city, endDate, startDate).Scan(&cars)

	result := m.App.DB.
		Table("cars c").
		Select("c.ID, c.model_id, c.type_id, c.location_id, c.Plate, c.Price, c.Color, c.Bags, c.Passengers, c.Year").
		Joins("LEFT JOIN car_histories ch ON c.ID = ch.Car_Id").
		Joins("LEFT JOIN rent_infos r ON ch.rent_info_id = r.ID").
		Joins("LEFT JOIN cities l ON c.location_id = l.ID").
		Where("l.name = ? AND (r.ID IS NULL OR (r.Start_Date > ? OR r.End_Date < ? OR r.status_id = 2 OR r.status_id = 3))", city, endDate, startDate).
		Order("c.ID").
		Preload("Type").
		Preload("Model").
		Preload("Location").
		Preload("Location.City").
		Preload("Model.Brand").
		Find(&cars)

	if result.Error != nil {
		m.App.Session.Put(r.Context(), "error", "Error finding available cars!")
		http.Redirect(w, r, "/cars", http.StatusSeeOther)
		return
	}

	var assignments []entities.CarAssignment
	if err = m.App.DB.Preload("Purpose").Preload("Car").Find(&assignments).Error; err != nil {
		http.Error(w, "Failed to fetch car assignments data", http.StatusInternalServerError)
		return
	}

	assignmentsByCarID := make(map[int][]entities.CarAssignment)
	for _, assignment := range assignments {
		carID := assignment.Car.ID
		assignmentsByCarID[carID] = append(assignmentsByCarID[carID], assignment)
	}

	for i, _ := range cars {
		cars[i].Assignments = assignmentsByCarID[cars[i].ID]
	}

	var cities []entities.City

	result = m.App.DB.Find(&cities)
	if result.Error != nil {
		http.Error(w, "Failed to fetch cities data", http.StatusInternalServerError)
		return
	}

	data := make(map[string]interface{})
	data["cars"] = cars
	data["Cities"] = cities

	rent = entities.RentInfo{
		StartDate: startDate,
		EndDate:   endDate,
	}
	m.App.Session.Put(r.Context(), "rent", rent)

	render.RenderTemplate(w, r, "cars.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) ChooseCar(w http.ResponseWriter, r *http.Request) {
	_, ok := m.App.Session.Get(r.Context(), "user_rent").(entities.RentInfo)
	if ok {
		m.ClearSessionData(r)
		m.App.Session.Put(r.Context(), "error", "Can't get user rent from session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	exploded := strings.Split(r.RequestURI, "/")
	carId, err := strconv.Atoi(exploded[2])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "missing url parameter")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	rent, ok := m.App.Session.Get(r.Context(), "rent").(entities.RentInfo)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	userId, ok := m.App.Session.Get(r.Context(), "user_id").(int)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Can't get user id from session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	rent.StatusId = 4
	rent.FromId = 1
	rent.ReturnId = 1
	rent.PaymentMethod = "none"
	res := m.App.DB.Create(&rent)
	if res.Error != nil {
		m.App.Session.Put(r.Context(), "error", "Can't create rent info")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var carRent entities.CarHistory
	carRent.CarId = carId
	carRent.RentInfoId = rent.ID
	res = m.App.DB.Create(&carRent)

	userHistory := entities.UserHistory{UserID: userId, RentInfoID: rent.ID}
	carHistory := entities.CarHistory{CarId: carId, RentInfoId: rent.ID}
	m.App.Session.Put(r.Context(), "rent", rent)
	m.App.Session.Put(r.Context(), "user_rent", userHistory)
	m.App.Session.Put(r.Context(), "car_rent", carHistory)

	http.Redirect(w, r, "/make-booking", http.StatusSeeOther)
}
