package controllers

import (
	models "help/models/app_models"
	"help/models/entities"
	"net/http"
)

var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *models.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *models.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewControllers sets the repository for the handlers
func NewControllers(r *Repository) {
	Repo = r
}

func (m *Repository) ClearSessionData(r *http.Request) {
	_, ok := m.App.Session.Get(r.Context(), "rent").(entities.RentInfo)
	if ok {
		m.App.Session.Remove(r.Context(), "rent")
	}
	_, ok = m.App.Session.Get(r.Context(), "car_rent").(entities.CarHistory)
	if ok {
		m.App.Session.Remove(r.Context(), "car_rent")
	}
	_, ok = m.App.Session.Get(r.Context(), "user_rent").(entities.UserHistory)
	if ok {
		m.App.Session.Remove(r.Context(), "user_rent")
	}
}
