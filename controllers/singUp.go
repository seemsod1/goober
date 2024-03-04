package controllers

import (
	"help/helpers/forms"
	"help/helpers/render"
	models "help/models/app_models"
	"help/models/entities"
	"net/http"
)

func (m *Repository) SingUpPage(w http.ResponseWriter, r *http.Request) {
	_, ok := m.App.Session.Get(r.Context(), "user_id").(int)
	if ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	render.RenderTemplate(w, r, "singUp.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) UserSingUp(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Can't parse form", http.StatusBadRequest)
		return
	}
	email := r.Form.Get("email")
	pass := r.Form.Get("password")
	name := r.Form.Get("name")
	phone := r.Form.Get("phone")

	form := forms.New(r.PostForm)
	form.Required("email", "password", "name", "phone")
	form.IsEmail("email")
	form.IsPhone("phone")
	form.MinLength("email", 5)
	form.Maxlength("email", 100)
	form.MinLength("password", 3)
	form.Maxlength("password", 100)
	form.MinLength("name", 1)
	form.Maxlength("name", 60)
	form.IsName("name")

	if !form.Valid() {
		render.RenderTemplate(w, r, "signUp.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	var user entities.User
	//find user by email
	if err = m.App.DB.Table("users").Where("email = ? OR phone = ?", email, phone).Take(&user).Error; err == nil {
		m.App.Session.Put(r.Context(), "error", "User already exists!")
		http.Redirect(w, r, "/join/singUp", http.StatusSeeOther)
		return
	}
	//if user don't exists

	user = entities.User{
		Name:     name,
		Email:    email,
		Phone:    phone,
		Password: pass,
		RoleId:   3,
	}
	//create user
	if err = m.App.DB.Create(&user).Error; err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't create user")
		http.Redirect(w, r, "/join/singUp", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Successfully singed up!")
	http.Redirect(w, r, "/join/login", http.StatusSeeOther)

}
