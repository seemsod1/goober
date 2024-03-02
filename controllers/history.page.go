package controllers

import (
	"help/helpers/render"
	models "help/models/app_models"
	"help/models/entities"
	"math"
	"net/http"
	"strconv"
	"strings"
)

func (m *Repository) MyHistory(w http.ResponseWriter, r *http.Request) {
	userId, _ := m.App.Session.Get(r.Context(), "user_id").(int)

	// Pagination parameters
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	} // You can get this from the request URL or any other source
	perPage := 10 // Number of items per page

	var userRents []entities.UserHistory
	res := m.App.DB.Table("user_histories").Where("user_id = ?", userId).Find(&userRents)

	var rents []entities.RentInfo

	for _, rent := range userRents {
		var r entities.RentInfo
		res = m.App.DB.Preload("Status").Table("rent_infos").First(&r, rent.RentInfoID)
		if res.Error != nil {
			continue
		}
		rents = append(rents, r)
	}

	var carsHistories []entities.CarHistory

	offset := (page - 1) * perPage

	for i := offset; i < min(offset+perPage, len(rents)); i++ {
		rent := rents[i]

		var cH entities.CarHistory
		res = m.App.DB.Preload("RentInfo").Preload("RentInfo.Status").Preload("RentInfo.From").Preload("RentInfo.From.City").Preload("RentInfo.Return").Preload("RentInfo.Return.City").Preload("Car").Preload("Car.Model").Preload("Car.Model.Brand").Table("car_histories").Where("rent_info_id = ?", rent.ID).Find(&cH)
		if res.Error != nil {
			continue
		}

		carsHistories = append(carsHistories, cH)
	}
	totalItems := len(rents)
	totalPages := int(math.Ceil(float64(totalItems) / float64(perPage)))

	data := make(map[string]interface{})
	data["carsHistories"] = carsHistories
	data["pagination"] = models.PaginationData{
		CurrentPage: page,
		TotalPages:  totalPages,
		PrevPage:    page - 1,
		NextPage:    page + 1,
		HasPrev:     page > 1,
		HasNext:     page < totalPages,
		Pages:       generatePageNumbers(page, totalPages),
	}

	render.RenderTemplate(w, r, "my-history.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// Helper function to generate page numbers for pagination links
func generatePageNumbers(currentPage, totalPages int) []int {
	var pages []int

	// You can customize how many pages you want to show in the pagination bar
	maxPages := 5
	start := currentPage - maxPages/2
	end := currentPage + maxPages/2

	if start < 1 {
		start = 1
	}

	if end > totalPages {
		end = totalPages
	}

	for i := start; i <= end; i++ {
		pages = append(pages, i)
	}

	return pages
}

func (m *Repository) FinishRent(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")

	rentId, err := strconv.Atoi(exploded[2])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "missing url parameter")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var rent entities.RentInfo
	res := m.App.DB.Preload("Status").Table("rent_infos").First(&rent, rentId)
	if res.Error != nil {
		m.App.Session.Put(r.Context(), "error", "Failed to find rent info")
		http.Redirect(w, r, "/my-history", http.StatusSeeOther)
		return
	}

	if rent.StatusId != 1 {
		m.App.Session.Put(r.Context(), "error", "Invalid rent status")
		http.Redirect(w, r, "/my-history", http.StatusSeeOther)
		return
	}

	rent.StatusId = 3
	rent.Status.ID = 3
	rent.Status.Name = "Finished"
	res = m.App.DB.Table("rent_infos").Save(&rent)
	if res.Error != nil {
		m.App.Session.Put(r.Context(), "error", "Failed to update rent status")
		http.Redirect(w, r, "/my-history", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Rent finished successfully")
	http.Redirect(w, r, "/my-history", http.StatusSeeOther)
}
