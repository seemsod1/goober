package controllers

import (
	"help/models/entities"
	"log"
	"net/http"
)

func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_, ok := m.App.Session.Get(r.Context(), "user").(entities.User)
	if !ok {
		log.Println("can't get item from session")
		m.App.Session.Put(r.Context(), "error", "Can't get user from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//m.App.Session.Put(r.Context(), "flash", "You've been logged out successfully!")
	//http.Redirect(w, r, "/join/login", http.StatusTemporaryRedirect)
}
