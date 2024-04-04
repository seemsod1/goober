package controllers

import (
	"github.com/google/uuid"
	"github.com/seemsod1/goober/helpers/render"
	models "github.com/seemsod1/goober/models/app_models"
	"github.com/seemsod1/goober/models/entities"
	"net/http"
)

func (m *Repository) ConfirmBooking(w http.ResponseWriter, r *http.Request) {
	job, ok := m.App.Session.Get(r.Context(), "task").(uuid.UUID)
	if !ok {
		m.ClearSessionData(r)
		m.App.Session.Put(r.Context(), "error", "Timeout! Please try again")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	flag := false
	jobs := m.App.Scheduler.Jobs()
	for _, j := range jobs {
		if j.ID() == job {
			flag = true
		}
	}
	if !flag {
		m.ClearSessionData(r)
		m.App.Session.Put(r.Context(), "error", "Timeout! Please try again")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	rent, ok := m.App.Session.Get(r.Context(), "rent").(entities.RentInfo)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	carRent, ok := m.App.Session.Get(r.Context(), "car_rent").(entities.CarHistory)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get car rent from session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var car entities.Car
	res := m.App.DB.Table("cars").Preload("Type").
		Preload("Model").
		Preload("Location").
		Preload("Location.City").
		Preload("Model.Brand").First(&car, carRent.CarId)

	if res.Error != nil {
		m.App.Session.Put(r.Context(), "error", "can't find car!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["Rent"] = rent
	data["Car"] = car

	render.RenderTemplate(w, r, "confirm.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) ConfirmBookingPost(w http.ResponseWriter, r *http.Request) {

	job, ok := m.App.Session.Get(r.Context(), "task").(uuid.UUID)
	if !ok {
		m.ClearSessionData(r)
		m.App.Session.Put(r.Context(), "error", "Timeout! Please try again")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	flag := false
	jobs := m.App.Scheduler.Jobs()
	for _, j := range jobs {
		if j.ID() == job {
			flag = true
		}
	}
	if !flag {
		m.ClearSessionData(r)
		m.App.Session.Put(r.Context(), "error", "Timeout! Please try again")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	rent, ok := m.App.Session.Get(r.Context(), "rent").(entities.RentInfo)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	rent.StatusId = 1
	m.App.DB.Save(&rent)
	err := m.App.Scheduler.RemoveJob(job)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Failed to find rent info")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var carRent entities.CarHistory
	res := m.App.DB.Preload("RentInfo").Preload("RentInfo.Return").Preload("RentInfo.Return.City").Preload("RentInfo.Status").Preload("Car").Preload("Car.Location").Preload("Car.Location.City").Table("car_histories").Where("rent_info_id = ?", rent.ID).First(&carRent)
	if res.Error != nil {
		m.App.Session.Put(r.Context(), "error", "Failed to find rent info")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//change in car location
	res = m.App.DB.Table("cars").Where("id = ?", carRent.CarId).Update("location_id", carRent.RentInfo.Return.ID)
	if res.Error != nil {
		m.App.Session.Put(r.Context(), "error", "Failed to update rent status")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")
	m.App.Session.Remove(r.Context(), "car_rent")
	m.App.Session.Remove(r.Context(), "rent")
	m.App.Session.Put(r.Context(), "flash", "Successfully booked!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
