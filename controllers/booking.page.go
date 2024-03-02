package controllers

import (
	"help/helpers/forms"
	"help/helpers/render"
	models "help/models/app_models"
	"help/models/entities"
	"net/http"
	"strconv"
	"time"
)

func (m *Repository) MakeBooking(w http.ResponseWriter, r *http.Request) {
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

	rent.Price = car.Price
	rent.From = car.Location
	rent.FromId = car.LocationId

	var locations []entities.RentLocation

	res = m.App.DB.Preload("City").Table("rent_locations").Not("id", car.LocationId).Find(&locations)
	if res.Error != nil {
		m.App.Session.Put(r.Context(), "error", "can't find locations!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	m.App.Session.Put(r.Context(), "rent", rent)
	time.AfterFunc(1*time.Minute, func() {
		_ = m.App.DB.Model(&entities.RentInfo{}).
			Where("status_id = ? AND created_at < ? AND id = ?", 4, time.Now().Add(-1*time.Minute), rent.ID).
			Updates(map[string]interface{}{"status_id": 2})
		re, _ := m.App.Session.Get(r.Context(), "rent").(entities.RentInfo)
		re.StatusId = 2
		m.App.Session.Put(r.Context(), "rent", re)
	})

	sd := rent.StartDate.Format("2006-01-02")
	ed := rent.EndDate.Format("2006-01-02")

	intmap := make(map[string]int)
	intmap["car_id"] = car.ID

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["locations"] = locations
	render.RenderTemplate(w, r, "booking.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		StringMap: stringMap,
		IntMap:    intmap,
		Data:      data,
	})
}

func (m *Repository) MakeBookingPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Can't parse form", http.StatusBadRequest)
		return
	}
	returnLocation := r.Form.Get("return_location")
	paymentMethod := r.Form.Get("payment_method")
	rent, ok := m.App.Session.Get(r.Context(), "rent").(entities.RentInfo)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	rLocid, err := strconv.Atoi(returnLocation)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse return location")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	var carLocation entities.RentLocation
	res := m.App.DB.Preload("City").First(&carLocation, rLocid)
	if res.Error != nil {
		m.App.Session.Put(r.Context(), "error", "can't find return location")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	rent.ReturnId = rLocid
	rent.Return = carLocation
	rent.PaymentMethod = paymentMethod
	m.App.Session.Put(r.Context(), "rent", rent)

	http.Redirect(w, r, "/confirm-booking", http.StatusSeeOther)
}