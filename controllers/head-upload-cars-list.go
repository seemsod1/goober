package controllers

import (
	"fmt"
	"github.com/seemsod1/goober/helpers/render"
	models "github.com/seemsod1/goober/models/app_models"
	"github.com/seemsod1/goober/models/entities"
	"github.com/xuri/excelize/v2"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

func (m *Repository) HeadUploadCarsList(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "head-upload-cars-list.page.tmpl", &models.TemplateData{})
}

func (m *Repository) HeadUploadCarsListPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Can't parse form", http.StatusBadRequest)
		m.App.Session.Put(r.Context(), "error", "Can't parse form")
		return
	}

	userId := m.App.Session.GetInt(r.Context(), "user_id")

	var location entities.RentLocation
	result := m.App.DB.Table("rent_locations").Where("user_id = ?", userId).Take(&location)
	if result.Error != nil {
		m.App.Session.Put(r.Context(), "error", "Location not found")
		return
	}

	file, _, err := r.FormFile("listUpload")
	if err != nil {
		http.Error(w, "Can't get file", http.StatusBadRequest)
		m.App.Session.Put(r.Context(), "error", "Can't get file")
		return
	}
	defer file.Close()

	tempFile, err := ioutil.TempFile("", "uploaded-cars-")
	if err != nil {
		http.Error(w, "Error creating temp file", http.StatusInternalServerError)
		return
	}
	defer func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		http.Error(w, "Error copying file", http.StatusInternalServerError)
		return
	}

	f, err := excelize.OpenFile(tempFile.Name())
	if err != nil {
		http.Error(w, "Can't read Excel file", http.StatusBadRequest)
		return
	}
	defer f.Close()

	rows, err := f.GetRows("cars")
	if err != nil {
		http.Error(w, "Can't get rows", http.StatusInternalServerError)
		return
	}

	var plates []string
	var col int
	var cars []entities.Car
	if len(rows) == 0 {
		m.App.Session.Put(r.Context(), "error", "Empty file")
		http.Redirect(w, r, "/head/upload-cars-list", http.StatusSeeOther)
		return
	}
	for i, row := range rows {
		var newCar bool
		newCar = false
		if i == 0 {
			col = len(row)
			continue
		}
		if len(row) == 0 {
			continue
		}

		if len(row)+1 != col {
			m.App.Session.Put(r.Context(), "error", "Error at row "+strconv.Itoa(i+1)+", invalid number of columns")
			http.Redirect(w, r, "/head/upload-cars-list", http.StatusSeeOther)
			return
		}
		var car entities.Car
		if row[0] == "" {
			m.App.Session.Put(r.Context(), "error", "Error at row "+strconv.Itoa(i+1)+", car type can't be empty")
			http.Redirect(w, r, "/head/upload-cars-list", http.StatusSeeOther)
			return
		}
		exist := m.App.DB.Where("name = ?", row[0]).First(&car.Type)
		if exist.Error != nil {
			m.App.Session.Put(r.Context(), "error", "Error at row "+strconv.Itoa(i+1)+", car type "+row[0]+" doesn't exist")
			http.Redirect(w, r, "/head/upload-cars-list", http.StatusSeeOther)
			return
		}
		if row[1] == "" {
			m.App.Session.Put(r.Context(), "error", "Error at row "+strconv.Itoa(i+1)+", brand can't be empty")
			http.Redirect(w, r, "/head/upload-cars-list", http.StatusSeeOther)
			return
		}
		exist = m.App.DB.Where("name = ?", row[1]).First(&car.Model.Brand)
		if exist.Error != nil {
			m.App.Session.Put(r.Context(), "error", "Error at row "+strconv.Itoa(i+1)+", car brand "+row[1]+" doesn't exist")
			http.Redirect(w, r, "/head/upload-cars-list", http.StatusSeeOther)
			return
		}
		model := row[2]
		if model == "" {
			m.App.Session.Put(r.Context(), "error", "Error at row "+strconv.Itoa(i+1)+", model can't be empty")
			http.Redirect(w, r, "/head/upload-cars-list", http.StatusSeeOther)
			return
		}
		exist = m.App.DB.Where("name = ? AND brand_id = ?", model, car.Model.Brand.ID).First(&car.Model)
		if exist.Error != nil {
			newCar = true
		}
		bags, err := strconv.Atoi(row[3])
		if err != nil || bags < 1 || bags > 6 {
			m.App.Session.Put(r.Context(), "error", "Error at row "+strconv.Itoa(i+1)+", bags must be a number")
			http.Redirect(w, r, "/head/upload-cars-list", http.StatusSeeOther)
			return
		}
		car.Bags = bags
		passengers, err := strconv.Atoi(row[4])
		if err != nil || passengers < 2 || passengers > 6 {
			m.App.Session.Put(r.Context(), "error", "Error at row "+strconv.Itoa(i+1)+", passengers must be a number")
			http.Redirect(w, r, "/head/upload-cars-list", http.StatusSeeOther)
			return
		}
		car.Passengers = passengers
		year, err := strconv.Atoi(row[5])
		if err != nil || year < 1950 || year > time.Now().Year() {
			m.App.Session.Put(r.Context(), "error", "Error at row "+strconv.Itoa(i+1)+", year must be a number")
			http.Redirect(w, r, "/head/upload-cars-list", http.StatusSeeOther)
			return
		}
		car.Year = year
		plate := row[6]
		if match, _ := regexp.MatchString("^[ABEIKMHOPCTXYZ]{2}\\d{4}[ABEIKMHOPCTXYZ]{2}$", plate); !match {
			m.App.Session.Put(r.Context(), "error", "Error at row "+strconv.Itoa(i+1)+", invalid plate number")
			http.Redirect(w, r, "/head/upload-cars-list", http.StatusSeeOther)
			return
		}
		for _, p := range plates {
			if p == plate {
				m.App.Session.Put(r.Context(), "error", "Error at row "+strconv.Itoa(i+1)+", plate number already exists")
				http.Redirect(w, r, "/head/upload-cars-list", http.StatusSeeOther)
				return
			}
		}
		var carByPlate entities.Car
		exist = m.App.DB.Where("plate = ?", plate).First(&carByPlate)
		if exist.Error == nil {
			car.ID = carByPlate.ID
			var rents []entities.CarHistory
			m.App.DB.Preload("RentInfo").Preload("RentInfo.Status").Where("car_id = ?", carByPlate.ID).Find(&rents)
			if len(rents) > 0 {
				for _, rent := range rents {
					if rent.RentInfo.Status.Name != "Finished" && rent.RentInfo.Status.Name != "Canceled" {
						m.App.Session.Put(r.Context(), "error", "Error at row "+strconv.Itoa(i+1)+", car with plate number "+plate+" has active rent or reservation")
						http.Redirect(w, r, "/head/upload-cars-list", http.StatusSeeOther)
						return
					}
				}
			}
		}

		plates = append(plates, plate)
		car.Plate = plate
		price, err := strconv.ParseFloat(row[7], 64)
		if err != nil || price < 1 || price > 10000 {
			m.App.Session.Put(r.Context(), "error", "Error at row "+strconv.Itoa(i+1)+", price must be a number")
			http.Redirect(w, r, "/head/upload-cars-list", http.StatusSeeOther)
			return
		}
		car.Price = price
		color := row[8]
		if match, _ := regexp.MatchString("^[A-Za-zА-Яа-яґҐєЄїЇ]+(\\s[A-Za-zА-Яа-я]+)*$", color); !match {
			m.App.Session.Put(r.Context(), "error", "Error at row "+strconv.Itoa(i+1)+", invalid color")
			http.Redirect(w, r, "/head/upload-cars-list", http.StatusSeeOther)
			return
		}
		if color == "" {
			m.App.Session.Put(r.Context(), "error", "Error at row "+strconv.Itoa(i+1)+", color can't be empty")
			http.Redirect(w, r, "/head/upload-cars-list", http.StatusSeeOther)
			return
		}
		//photo
		axis, err := excelize.CoordinatesToCellName(col, i+1)
		if err != nil {
			http.Error(w, "Can't get cell coordinate", http.StatusInternalServerError)
			return
		}
		pics, err := f.GetPictures("cars", axis)

		if err != nil {
			http.Error(w, "Can't get pictures", http.StatusInternalServerError)
			return
		}

		if len(pics) == 0 {
			m.App.Session.Put(r.Context(), "error", "Error at row "+strconv.Itoa(i+1)+", photo can't be empty")
			http.Redirect(w, r, "/head/upload-cars-list", http.StatusSeeOther)
			return
		}

		newFileName := fmt.Sprintf("%s%s", model, ".png")

		filePath := filepath.Join("D:\\golang\\car-rent\\resources\\img\\cars", newFileName)

		if _, err := os.Stat(filePath); err == nil {
			if err := os.Remove(filePath); err != nil {
				fmt.Println("Error deleting file:", err)
			}
		}

		if err := os.WriteFile(filePath, pics[0].File, 0644); err != nil {
			fmt.Println("Error saving file:", err)
		}

		car.Color = color
		car.LocationId = location.ID
		car.TypeId = car.Type.ID
		if newCar {
			var carModel entities.CarModel
			carModel.BrandId = car.Model.Brand.ID
			carModel.Name = model
			if err := m.App.DB.Create(&carModel).Error; err != nil {
				m.App.Session.Put(r.Context(), "error", "Error at row "+strconv.Itoa(i+1)+", failed to save car model")
				http.Redirect(w, r, "/head/upload-cars-list", http.StatusSeeOther)
				return
			}
			car.ModelId = carModel.ID
			car.Model = carModel
		}
		cars = append(cars, car)
		i++

	}

	for _, car := range cars {
		m.App.DB.Save(&car)
	}

	m.App.Session.Put(r.Context(), "flash", "Successfully added cars!")
	http.Redirect(w, r, "/head/", http.StatusSeeOther)

}
