package controllers

import (
	"github.com/seemsod1/goober/helpers/render"
	models "github.com/seemsod1/goober/models/app_models"
	"github.com/seemsod1/goober/models/entities"
	"net/http"
)

func (m *Repository) AdminPage(w http.ResponseWriter, r *http.Request) {

	// rents
	var rents []entities.RentInfo
	res := m.App.DB.Order("created_at desc").Find(&rents)

	var rentlen = len(rents)

	//revenue
	var revenue float64
	for _, rent := range rents {
		if rent.StatusId == 1 || rent.StatusId == 3 {
			revenue += rent.Price
		}
	}

	//cancelled
	var cancelled int
	for _, rent := range rents {
		if rent.StatusId == 2 {
			cancelled++
		}
	}

	//cars
	var clients int64
	res = m.App.DB.Table("users").Where("role_id = ?", 3).Count(&clients)

	data := make(map[string]interface{})
	data["Revenue"] = revenue

	var carHistories []entities.CarHistory
	for _, rent := range rents {
		var cH entities.CarHistory
		res = m.App.DB.Preload("Car").Preload("Car.Model").Preload("Car.Model.Brand").Table("car_histories").Where("rent_info_id = ?", rent.ID).Find(&cH)
		if res.Error != nil {
			continue
		}
		carHistories = append(carHistories, cH)
	}

	brandNum := make(map[string]int)
	for _, carHistory := range carHistories {
		brand := carHistory.Car.Model.Brand.Name
		brandNum[brand]++
	}
	receivedMap := make(map[string]int)
	count := 0
	for key, value := range brandNum {
		receivedMap[key] = value
		count++
		if count == 5 {
			break
		}
	}

	data["BrandNum"] = receivedMap

	intMap := make(map[string]int)
	intMap["Rents"] = rentlen
	intMap["Clients"] = int(clients)
	intMap["Cancelled"] = cancelled

	render.RenderTemplate(w, r, "admin.page.tmpl", &models.TemplateData{
		Data:   data,
		IntMap: intMap,
	})
}
