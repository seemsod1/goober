package main

import (
	"github.com/go-chi/chi/v5"
	"help/controllers"
	"net/http"
)

func routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(SessionLoad)
	mux.NotFound(controllers.Repo.NotFoundPage)

	mux.Group(func(mux chi.Router) {
		mux.Use(NoSurf)
		mux.Get("/", controllers.Repo.HomePage)
		mux.Get("/about", controllers.Repo.About)
		mux.Get("/cars", controllers.Repo.CarsPage)
		mux.Post("/cars", controllers.Repo.CarsPagePost)
	})
	mux.Group(func(mux chi.Router) {
		mux.Use(RequireAuth)
		mux.Use(User)
		mux.Get("/choose-car/{id}", controllers.Repo.ChooseCar)
		mux.Get("/make-booking", controllers.Repo.MakeBooking)
		mux.Post("/make-booking", controllers.Repo.MakeBookingPost)
		mux.Get("/confirm-booking", controllers.Repo.ConfirmBooking)
		mux.Post("/confirm-booking", controllers.Repo.ConfirmBookingPost)
		mux.Get("/my-history", controllers.Repo.MyHistory)
		mux.Post("/finish-rent/{id}", controllers.Repo.FinishRent)
	})

	join := chi.NewRouter()
	join.Use(SessionLoad)
	join.Group(func(join chi.Router) {
		join.Use(NoSurf)
		join.Get("/login", controllers.Repo.LoginPage)
		join.Post("/login", controllers.Repo.UserLogin)
	})
	join.Get("/singUp", controllers.Repo.SingUpPage)
	join.Post("/singUp", controllers.Repo.UserSingUp)

	join.Get("/logout", controllers.Repo.Logout)
	mux.Mount("/join", join)

	head := chi.NewRouter()
	head.Use(SessionLoad)
	head.Use(Head)
	head.Use(RequireAuth)
	head.Use(NoSurf)

	head.Get("/car-history/{id}", controllers.Repo.CarHistory)
	head.Get("/all-cars", controllers.Repo.AllCars)
	head.Get("/rents-histories", controllers.Repo.RentsHistories)
	head.Get("/add-car-model", controllers.Repo.AddCarModel)
	head.Post("/add-car-model", controllers.Repo.AddCarModelPost)
	head.Get("/", controllers.Repo.HeadPage)
	head.Get("/add-car", controllers.Repo.AddCar)
	head.Get("/get-brands-with-types", controllers.Repo.GetBrandsWithTypes)
	head.Get("/get-brands", controllers.Repo.GetAllBrands)
	head.Get("/get-models", controllers.Repo.GetModels)
	head.Get("/get-types", controllers.Repo.GetTypes)
	head.Post("/add-car", controllers.Repo.AddCarPost)
	head.Post("/change-car-price", controllers.Repo.ChangeCarPrice)
	head.Get("/save-cars-list", controllers.Repo.SaveCarsList)
	head.Post("/save-cars-list", controllers.Repo.SaveCarsListPost)
	head.Get("/upload-cars-list", controllers.Repo.UploadCarsList)
	head.Post("/upload-cars-list", controllers.Repo.UploadCarsListPost)

	mux.Mount("/head", head)

	fileServer := http.FileServer(http.Dir("./resources/"))
	mux.Handle("/resources/*", http.StripPrefix("/resources", fileServer))
	return mux
}
