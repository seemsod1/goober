package controllers

import (
	"github.com/seemsod1/goober/helpers"
	"github.com/seemsod1/goober/helpers/forms"
	"github.com/seemsod1/goober/helpers/render"
	models "github.com/seemsod1/goober/models/app_models"
	"github.com/seemsod1/goober/models/entities"
	"net/http"
)

func (m *Repository) ChangePassword(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "change-password.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) ChangePasswordPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't parse form")
		http.Redirect(w, r, "/change-password", http.StatusSeeOther)
		return
	}
	oldPass := r.Form.Get("old_pass")
	newPass := r.Form.Get("new_pass")
	newPassRepeat := r.Form.Get("new_pass_repeat")

	form := forms.New(r.PostForm)
	form.Required("old_pass", "new_pass", "new_pass_repeat")
	form.MinLength("new_pass", 3)
	form.Maxlength("new_pass", 100)
	form.MinLength("old_pass", 3)
	form.Maxlength("old_pass", 100)
	form.MinLength("new_pass_repeat", 3)
	form.Maxlength("new_pass_repeat", 100)
	form.IsPasswordMatch("new-password", "confirm-password")
	if !form.Valid() {
		render.RenderTemplate(w, r, "change-password.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	userId, ok := m.App.Session.Get(r.Context(), "user_id").(int)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Can't get user id")
		http.Redirect(w, r, "/change-password", http.StatusSeeOther)
		return
	}

	var user entities.User
	if err = m.App.DB.Table("users").Where("id = ?", userId).Take(&user).Error; err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't find user")
		http.Redirect(w, r, "/change-password", http.StatusSeeOther)
		return
	}

	if oldPass != user.Password {
		m.App.Session.Put(r.Context(), "error", "Wrong old password")
		http.Redirect(w, r, "/change-password", http.StatusSeeOther)
		return
	}

	if !helpers.CheckPassword(newPass, newPassRepeat) {
		m.App.Session.Put(r.Context(), "error", "Passwords don't match")
		http.Redirect(w, r, "/change-password", http.StatusSeeOther)
		return
	}

	user.Password = newPass
	m.App.DB.Save(&user)

	m.App.Session.Put(r.Context(), "flash", "Password changed")
	http.Redirect(w, r, "/change-password", http.StatusSeeOther)
}
