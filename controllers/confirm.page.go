package controllers

import (
	"help/helpers/render"
	models "help/models/app_models"
	"help/models/entities"
	"net/http"
	"time"
)

func (m *Repository) ConfirmBooking(w http.ResponseWriter, r *http.Request) {
	rent, ok := m.App.Session.Get(r.Context(), "rent").(entities.RentInfo)
	if !ok || rent.StatusId != 4 {
		m.ClearSessionData(r)
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

	time.AfterFunc(1*time.Minute, func() {
		m.App.Session.Remove(r.Context(), "rent")
		m.App.Session.Remove(r.Context(), "user_rent")
		m.App.Session.Remove(r.Context(), "car_rent")
		_ = m.App.DB.Model(&entities.RentInfo{}).
			Where("status_id = ? AND created_at < ? AND id = ?", 4, time.Now().Add(-1*time.Minute), rent.ID).
			Updates(map[string]interface{}{"status_id": 2})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	render.RenderTemplate(w, r, "confirm.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) ConfirmBookingPost(w http.ResponseWriter, r *http.Request) {
	rent, ok := m.App.Session.Get(r.Context(), "rent").(entities.RentInfo)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	rent.StatusId = 1
	m.App.DB.Save(&rent)

	m.App.Session.Remove(r.Context(), "reservation")
	m.App.Session.Remove(r.Context(), "car_rent")
	m.App.Session.Remove(r.Context(), "rent")
	m.App.Session.Put(r.Context(), "flash", "Successfully booked!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
