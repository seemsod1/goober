package controllers

import (
	"encoding/json"
	"fmt"
	"help/helpers/forms"
	"help/helpers/render"
	models "help/models/app_models"
	"help/models/entities"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func (m *Repository) HeadAddCarModel(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "head-add-car-model.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) HeadAddCarModelPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Can't parse form", http.StatusBadRequest)
		return
	}

	brand := r.Form.Get("SelectBrand")
	model := r.Form.Get("inputModel")

	form := forms.New(r.PostForm)
	form.Required("SelectBrand", "inputModel")
	form.IsNumber("SelectBrand")
	form.MinLength("inputModel", 1)
	form.Maxlength("inputModel", 60)

	if !form.Valid() {
		render.RenderTemplate(w, r, "head-add-car-model.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}
	brandId, _ := strconv.Atoi(brand)

	var carModel entities.CarModel

	result := m.App.DB.Where("LOWER(name) = LOWER(?) AND LOWER(brand_id) = LOWER(?)", model, brandId).First(&carModel)
	if result.Error == nil {
		m.App.Session.Put(r.Context(), "error", "Car model already exists")
		http.Redirect(w, r, "/head/", http.StatusSeeOther)
		return
	}

	carModel.BrandId = brandId
	carModel.Name = model

	result = m.App.DB.Create(&carModel)
	if result.Error != nil {
		http.Error(w, "Failed to save car model", http.StatusInternalServerError)
		return
	}

	file, handler, err := r.FormFile("photoUpload")
	if err == nil {
		defer file.Close()

		if strings.ToLower(filepath.Ext(handler.Filename)) != ".png" {
			m.App.Session.Put(r.Context(), "error", "Only PNG!")
			return
		}

		newFileName := fmt.Sprintf("%s%s", model, ".png")

		filePath := filepath.Join("D:\\golang\\car-rent\\resources\\img\\cars", newFileName)

		dst, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Failed to save the file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Failed to copy the file", http.StatusInternalServerError)
			return
		}
	}

	m.App.Session.Put(r.Context(), "flash", "Successfully add car model!")
	http.Redirect(w, r, "/head/", http.StatusSeeOther)

}

func (m *Repository) HeadGetAllBrands(w http.ResponseWriter, r *http.Request) {
	var brands []entities.CarBrand

	result := m.App.DB.Find(&brands)
	if result.Error != nil {
		http.Error(w, "Failed to fetch brands data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(brands)
}
