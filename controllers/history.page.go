package controllers

import (
	"encoding/json"
	"help/helpers"
	"help/helpers/render"
	models "help/models/app_models"
	"help/models/entities"
	"log"
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
	res := m.App.DB.Table("user_histories").Where("user_id = ?", userId).Order("created_at desc").Find(&userRents)

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
		Pages:       helpers.GeneratePageNumbers(page, totalPages),
	}

	render.RenderTemplate(w, r, "my-history.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func (m *Repository) FinishRent(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")

	rentId, err := strconv.Atoi(exploded[2])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "missing url parameter")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var carRent entities.CarHistory
	res := m.App.DB.Preload("RentInfo").Preload("RentInfo.Return").Preload("RentInfo.Return.City").Preload("RentInfo.Status").Preload("Car").Preload("Car.Location").Preload("Car.Location.City").Table("car_histories").Where("rent_info_id = ?", rentId).First(&carRent)
	if res.Error != nil {
		m.App.Session.Put(r.Context(), "error", "Failed to find rent info")
		http.Redirect(w, r, "/my-history", http.StatusSeeOther)
		return
	}

	if carRent.RentInfo.StatusId != 1 {
		m.App.Session.Put(r.Context(), "error", "Invalid rent status")
		http.Redirect(w, r, "/my-history", http.StatusSeeOther)
		return
	}

	res = m.App.DB.Table("rent_infos").Where("id = ?", rentId).Update("status_id", 3)
	if res.Error != nil {
		m.App.Session.Put(r.Context(), "error", "Failed to update rent status")
		http.Redirect(w, r, "/my-history", http.StatusSeeOther)
		return
	}

	resp := jsonResponse{
		OK:      true,
		Message: "Successful finished rent!",
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
