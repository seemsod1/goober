package controllers

import (
	"net/http"
)

func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_, ok := m.App.Session.Get(r.Context(), "user_id").(int)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	m.App.Session.Put(r.Context(), "flash", "Logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
