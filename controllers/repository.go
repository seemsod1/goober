package controllers

import models "help/models/app_models"

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
