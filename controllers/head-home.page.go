package controllers

import (
	"help/helpers/render"
	models "help/models/app_models"
	"help/models/entities"
	"net/http"
)

func (m *Repository) HeadPage(w http.ResponseWriter, r *http.Request) {
	userId, _ := m.App.Session.Get(r.Context(), "user_id").(int)

	var location entities.RentLocation
	res := m.App.DB.Table("rent_locations").Where("user_id = ?", userId).Take(&location)
	if res.Error != nil {
		m.App.Session.Put(r.Context(), "error", "Can't get locations")
		return
	}
	// rents
	var rents []entities.RentInfo
	res = m.App.DB.Table("rent_infos").Where("from_id = ?", location.ID).Order("created_at desc").Find(&rents)

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
	var cars int64
	res = m.App.DB.Table("cars").Where("location_id = ?", location.ID).Count(&cars)

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
	intMap["Cars"] = int(cars)
	intMap["Cancelled"] = cancelled

	render.RenderTemplate(w, r, "head.page.tmpl", &models.TemplateData{
		Data:   data,
		IntMap: intMap,
	})
}
