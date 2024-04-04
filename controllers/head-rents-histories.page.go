package controllers

import (
	"github.com/seemsod1/goober/helpers"
	"github.com/seemsod1/goober/helpers/render"
	models "github.com/seemsod1/goober/models/app_models"
	"github.com/seemsod1/goober/models/entities"
	"math"
	"net/http"
	"strconv"
)

func (m *Repository) HeadRentsHistories(w http.ResponseWriter, r *http.Request) {
	userId, _ := m.App.Session.Get(r.Context(), "user_id").(int)

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	perPage := 10

	var location entities.RentLocation
	res := m.App.DB.Table("rent_locations").Where("user_id = ?", userId).Take(&location)
	if res.Error != nil {
		m.App.Session.Put(r.Context(), "error", "Can't get locations")
		return
	}

	var rents []entities.RentInfo
	res = m.App.DB.Table("rent_infos").Where("from_id = ?", location.ID).Order("created_at desc").Find(&rents)

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

	render.RenderTemplate(w, r, "head-rents-histories.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
