package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/seemsod1/goober/controllers"
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
		mux.Use(NoSurf)
		mux.Use(RequireAuth)
		mux.Get("/verify-mail", controllers.Repo.VerifyMail)
		mux.Post("/verify-mail", controllers.Repo.VerifyMailPost)
		mux.Post("/get-verification-code", controllers.Repo.GetVerificationCode)

		mux.Get("/change-password", controllers.Repo.ChangePassword)
		mux.Post("/change-password", controllers.Repo.ChangePasswordPost)

		mux.Post("/connect-metamask", controllers.Repo.PaymentsConnectMetamask)
	})

	mux.Group(func(mux chi.Router) {
		mux.Use(RequireAuth)
		mux.Use(RequireVerified)
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
	head.Use(RequireVerified)
	head.Use(NoSurf)

	head.Get("/car-history/{id}", controllers.Repo.HeadCarHistory)
	head.Get("/all-cars", controllers.Repo.HeadAllCars)
	head.Get("/rents-histories", controllers.Repo.HeadRentsHistories)
	head.Get("/add-car-model", controllers.Repo.HeadAddCarModel)
	head.Post("/add-car-model", controllers.Repo.HeadAddCarModelPost)
	head.Get("/", controllers.Repo.HeadPage)
	head.Get("/add-car", controllers.Repo.HeadAddCar)
	head.Get("/get-brands-with-types", controllers.Repo.HeadGetBrandsWithTypes)
	head.Get("/get-brands", controllers.Repo.HeadGetAllBrands)
	head.Get("/get-models", controllers.Repo.HeadGetModels)
	head.Get("/get-types", controllers.Repo.HeadGetTypes)
	head.Post("/add-car", controllers.Repo.HeadAddCarPost)
	head.Post("/change-car-price", controllers.Repo.HeadChangeCarPrice)
	head.Post("/change-car-photo", controllers.Repo.HeadChangeCarPhoto)
	head.Get("/save-cars-list", controllers.Repo.HeadSaveCarsList)
	head.Post("/save-cars-list", controllers.Repo.HeadSaveCarsListPost)
	head.Get("/upload-cars-list", controllers.Repo.HeadUploadCarsList)
	head.Post("/upload-cars-list", controllers.Repo.HeadUploadCarsListPost)

	mux.Mount("/head", head)

	admin := chi.NewRouter()
	admin.Use(SessionLoad)
	admin.Use(Admin)
	admin.Use(RequireAuth)
	admin.Use(RequireVerified)
	admin.Use(NoSurf)
	admin.Get("/", controllers.Repo.AdminPage)
	admin.Get("/all-users", controllers.Repo.AdminAllUsers)
	admin.Get("/delete-user/{id}", controllers.Repo.AdminDeleteUser)
	admin.Get("/locations", controllers.Repo.AdminGetLocations)
	admin.Post("/promote-user", controllers.Repo.AdminPromoteUser)

	mux.Mount("/admin", admin)

	fileServer := http.FileServer(http.Dir("./resources/"))
	mux.Handle("/resources/*", http.StripPrefix("/resources", fileServer))
	return mux
}
