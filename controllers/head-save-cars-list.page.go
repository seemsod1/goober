package controllers

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"help/helpers/render"
	models "help/models/app_models"
	"help/models/entities"
	"net/http"
)

func (m *Repository) HeadSaveCarsList(w http.ResponseWriter, r *http.Request) {
	userId := m.App.Session.GetInt(r.Context(), "user_id")

	var location entities.RentLocation
	result := m.App.DB.Table("rent_locations").Where("user_id = ?", userId).Take(&location)
	if result.Error != nil {
		m.App.Session.Put(r.Context(), "error", "Location not found")
		return
	}

	SheetName := "cars"
	f := excelize.NewFile()

	index, err := f.NewSheet("Sheet1")
	f.SetActiveSheet(index)
	f.SetSheetName("Sheet1", SheetName)

	err = f.SetColWidth(SheetName, "A", "A", 10) //type
	err = f.SetColWidth(SheetName, "B", "B", 10) // brand
	err = f.SetColWidth(SheetName, "C", "C", 20) // model
	err = f.SetColWidth(SheetName, "D", "D", 6)  // bags
	err = f.SetColWidth(SheetName, "E", "E", 15) // passengers
	err = f.SetColWidth(SheetName, "F", "F", 8)  // year
	err = f.SetColWidth(SheetName, "G", "G", 10) // plate
	err = f.SetColWidth(SheetName, "H", "H", 8)  // price
	err = f.SetColWidth(SheetName, "I", "I", 10) // color

	style, err := f.NewStyle(&excelize.Style{Font: &excelize.Font{Size: 14, Bold: true}})
	if err != nil {
		fmt.Println(err)
	}

	err = f.SetRowStyle(SheetName, 1, 1, style)

	if err != nil {
		fmt.Println(err)
	}

	var cars []entities.Car
	result = m.App.DB.Table("cars").Preload("Type").Preload("Model").Preload("Model.Brand").Where("location_id = ?", location.ID).Find(&cars)
	if result.Error != nil {
		m.App.Session.Put(r.Context(), "error", "Can't get cars")
		return
	}

	// Set header row
	headerRow := []string{"Type", "Brand", "Model", "Bags", "Passengers", "Year", "Plate", "Price", "Color", "", "", "Types", "Brands", "Models"}
	for colIdx, value := range headerRow {
		cellCoordinate := fmt.Sprintf("%s%d", string('A'+colIdx), 1)
		err = f.SetCellValue(SheetName, cellCoordinate, value)
		if err != nil {
			fmt.Println(err)
		}
	}

	// Set data rows
	for rowIdx, car := range cars {
		rowData := []interface{}{car.Type.Name, car.Model.Brand.Name, car.Model.Name, car.Bags, car.Passengers, car.Year, car.Plate, car.Price, car.Color}
		for colIdx, value := range rowData {
			cellCoordinate := fmt.Sprintf("%s%d", string('A'+colIdx), rowIdx+2)
			err = f.SetCellValue(SheetName, cellCoordinate, value)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	if err = f.SaveAs("D:\\golang\\car-rent\\resources\\excel\\cars.xlsx"); err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't get file")
		http.Error(w, "Can't parse form", http.StatusInternalServerError)
	}

	render.RenderTemplate(w, r, "head-save-cars-list.page.tmpl", &models.TemplateData{})
}

func (m *Repository) HeadSaveCarsListPost(w http.ResponseWriter, r *http.Request) {

	m.App.Session.Put(r.Context(), "flash", "Successfully saved!")
	http.Redirect(w, r, "/head/", http.StatusSeeOther)
}
