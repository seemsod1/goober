package controllers

import (
	"github.com/go-co-op/gocron/v2"
	"github.com/seemsod1/goober/helpers/render"
	models "github.com/seemsod1/goober/models/app_models"
	"github.com/seemsod1/goober/models/entities"
	"log"
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
		m.ClearSessionData(r)
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

	result := m.App.DB.
		Preload("Location").
		Preload("Location.City").
		Preload("Model").
		Preload("Model.Brand").
		Preload("Type").
		Table("cars").
		Select("cars.*").
		Joins("JOIN rent_locations ON cars.location_id = rent_locations.id").
		Joins("JOIN cities ON rent_locations.city_id = cities.id").
		Where("cities.name = ?", city).
		Not("cars.id IN (?)", m.App.DB.Table("car_histories").
			Select("DISTINCT car_id").
			Joins("JOIN rent_infos ON car_histories.rent_info_id = rent_infos.id").
			Where("rent_infos.status_id IN (?)", []int{4, 1}).
			Where(
				"(? BETWEEN DATE(rent_infos.start_date) AND DATE(rent_infos.end_date)) OR "+
					"(? BETWEEN DATE(rent_infos.start_date) AND DATE(rent_infos.end_date)) OR "+
					"(DATE(rent_infos.start_date) <= ? AND DATE(rent_infos.end_date) >= ?)", startDate, endDate, startDate, endDate),
		).
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

	var car entities.Car
	result := m.App.DB.First(&car, carId)
	if result.Error != nil {
		m.App.Session.Put(r.Context(), "error", "Can't get car from db")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	rent.StatusId = 4
	rent.FromId = car.LocationId
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
	if res.Error != nil {
		m.App.Session.Put(r.Context(), "error", "Can't create car history")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var userRent entities.UserHistory
	userRent.UserID = userId
	userRent.RentInfoID = rent.ID
	res = m.App.DB.Create(&userRent)
	if res.Error != nil {
		m.App.Session.Put(r.Context(), "error", "Can't create user history")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	j, err := m.App.Scheduler.NewJob(
		gocron.DurationJob(
			time.Minute*1,
		),
		gocron.NewTask(func() {
			_ = m.App.DB.Model(&entities.RentInfo{}).
				Where("status_id = ? AND id = ?", 4, rent.ID).
				Updates(map[string]interface{}{"status_id": 2})
			m.App.Session.Put(r.Context(), "error", "Time is up! Your reservation has been canceled!")
			log.Println("Task is done!")
		}),
		gocron.WithLimitedRuns(1),
	)
	m.App.Session.Put(r.Context(), "task", j.ID())
	log.Println("Job ID: ", j.ID())

	userHistory := entities.UserHistory{UserID: userId, RentInfoID: rent.ID}
	carHistory := entities.CarHistory{CarId: carId, RentInfoId: rent.ID}
	m.App.Session.Put(r.Context(), "rent", rent)
	m.App.Session.Put(r.Context(), "user_rent", userHistory)
	m.App.Session.Put(r.Context(), "car_rent", carHistory)

	http.Redirect(w, r, "/make-booking", http.StatusSeeOther)
}
