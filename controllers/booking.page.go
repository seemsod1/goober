package controllers

import (
	"github.com/google/uuid"
	"github.com/seemsod1/goober/helpers/forms"
	"github.com/seemsod1/goober/helpers/render"
	models "github.com/seemsod1/goober/models/app_models"
	"github.com/seemsod1/goober/models/entities"
	"net/http"
	"strconv"
)

func (m *Repository) MakeBooking(w http.ResponseWriter, r *http.Request) {
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
	if rent.EndDate == rent.StartDate {
		rent.Price = car.Price
	} else {
		difference := rent.EndDate.Sub(rent.StartDate).Hours() / 24
		rent.Price = car.Price * float64(difference)
	}
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
