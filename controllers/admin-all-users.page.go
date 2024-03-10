package controllers

import (
	"encoding/json"
	"help/helpers"
	"help/helpers/forms"
	"help/helpers/render"
	models "help/models/app_models"
	"help/models/entities"
	"math"
	"net/http"
	"strconv"
	"strings"
)

func (m *Repository) AdminAllUsers(w http.ResponseWriter, r *http.Request) {

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	perPage := 10

	var users []entities.User
	var totalItems int64
	res := m.App.DB.Table("users").Preload("Role").Order("role_id asc,updated_at").Where("role_id != 1").Offset((page - 1) * perPage).Limit(perPage).Find(&users).Count(&totalItems)

	if res.Error != nil {
		http.Error(w, res.Error.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(perPage)))

	data := make(map[string]interface{})
	data["users"] = users
	data["pagination"] = models.PaginationData{
		CurrentPage: page,
		TotalPages:  totalPages,
		PrevPage:    page - 1,
		NextPage:    page + 1,
		HasPrev:     page > 1,
		HasNext:     page < totalPages,
		Pages:       helpers.GeneratePageNumbers(page, totalPages),
	}

	render.RenderTemplate(w, r, "admin-all-users.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) AdminDeleteUser(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	userId, err := strconv.Atoi(exploded[3])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "missing url parameter")
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	var canc entities.RentStatus
	m.App.DB.Where("name = ?", "Cancelled").First(&canc)
	var userHistory []entities.UserHistory
	res := m.App.DB.Preload("RentInfo").Preload("RentInfo.Status").Where("user_id = ?", userId).Find(&userHistory)
	if res.Error == nil {
		for _, uh := range userHistory {
			if uh.RentInfo.Status.Name == "Active" || uh.RentInfo.Status.Name == "Processing" {
				uh.RentInfo.Status = canc
				uh.RentInfo.Status.ID = canc.ID

				m.App.DB.Save(&uh.RentInfo)
			}
		}
	}

	res = m.App.DB.Delete(&entities.User{}, userId)
	if res.Error != nil {
		m.App.Session.Put(r.Context(), "error", "Can't delete user")
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "User deleted")
	http.Redirect(w, r, "/admin/all-users", http.StatusSeeOther)
}

func (m *Repository) AdminGetLocations(w http.ResponseWriter, r *http.Request) {
	var locations []entities.RentLocation
	res := m.App.DB.Preload("City").Find(&locations)
	if res.Error != nil {
		http.Error(w, res.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(locations)
}

func (m *Repository) AdminPromoteUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Can't parse form", http.StatusBadRequest)
		return
	}
	loc := r.Form.Get("location")
	userForm := r.Form.Get("user")

	form := forms.New(r.PostForm)
	form.Required("location", "user")
	form.IsNumber("location")
	form.IsNumber("user")

	if !form.Valid() {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	locationId, err := strconv.Atoi(loc)
	if err != nil {
		http.Error(w, "Invalid location", http.StatusBadRequest)
		return
	}
	userId, err := strconv.Atoi(userForm)
	if err != nil {
		http.Error(w, "Invalid user", http.StatusBadRequest)
		return
	}

	var user entities.User
	res := m.App.DB.First(&user, userId)
	if res.Error != nil {
		m.App.Session.Put(r.Context(), "error", "Can't find user")
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	var managerRole entities.UserRole
	res = m.App.DB.Where("name = ?", "Manager").First(&managerRole)
	if res.Error != nil {
		m.App.Session.Put(r.Context(), "error", "Can't find role")
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	var userRole entities.UserRole
	res = m.App.DB.Where("name = ?", "User").First(&userRole)
	if res.Error != nil {
		m.App.Session.Put(r.Context(), "error", "Can't find role")
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return

	}
	var location entities.RentLocation
	res = m.App.DB.Preload("User").Preload("User.Role").First(&location, locationId)
	if res.Error != nil {
		m.App.Session.Put(r.Context(), "error", "Can't find location")
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	prevManager := location.User
	if prevManager.ID != 0 {
		prevManager.Role = userRole
		prevManager.RoleId = userRole.ID
		res = m.App.DB.Save(&prevManager)
		if res.Error != nil {
			m.App.Session.Put(r.Context(), "error", "Can't demote user")
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			return
		}
	}

	user.Role = managerRole
	user.RoleId = managerRole.ID
	res = m.App.DB.Save(&user)
	if res.Error != nil {
		m.App.Session.Put(r.Context(), "error", "Can't promote user")
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	location.User = user
	location.UserId = user.ID
	res = m.App.DB.Save(&location)
	if res.Error != nil {
		m.App.Session.Put(r.Context(), "error", "Can't save location")
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return

	}

	http.Redirect(w, r, "/admin/all-users", http.StatusSeeOther)
}
