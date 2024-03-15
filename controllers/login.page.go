package controllers

import (
	"help/helpers"
	"help/helpers/forms"
	"help/helpers/render"
	models "help/models/app_models"
	"help/models/entities"
	"net/http"
)

func (m *Repository) LoginPage(w http.ResponseWriter, r *http.Request) {
	_, ok := m.App.Session.Get(r.Context(), "user_id").(int)
	if ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	render.RenderTemplate(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})

}

func (m *Repository) UserLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Can't parse form", http.StatusBadRequest)
		return
	}
	email := r.Form.Get("email")
	pass := r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")
	form.MinLength("email", 5)
	form.Maxlength("email", 100)
	form.MinLength("password", 3)
	form.Maxlength("password", 100)
	form.IsEmail("email")

	if !form.Valid() {
		render.RenderTemplate(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	var user entities.User
	//find user by email
	if err = m.App.DB.Table("users").Where("email = ?", email).Take(&user).Error; err != nil {
		m.App.Session.Put(r.Context(), "error", "Cant find user!")
		http.Redirect(w, r, "/join/login", http.StatusSeeOther)
		return
	}
	//if user exists
	if !helpers.CheckPassword(user.Password, pass) {
		m.App.Session.Put(r.Context(), "error", "Invalid password!")
		http.Redirect(w, r, "/join/login", http.StatusSeeOther)
		return
	}
	m.App.Session.Put(r.Context(), "user_id", user.ID)
	m.App.Session.Put(r.Context(), "user_role", user.RoleId)
	if user.IsVerified == false {
		m.App.Session.Put(r.Context(), "verified", user.IsVerified)
	}
	m.App.Session.Put(r.Context(), "flash", "Successfully logged in!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
