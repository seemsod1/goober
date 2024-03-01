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
		mux.Get("/choose-car/{id}", controllers.Repo.ChooseCar)
		mux.Get("/make-booking", controllers.Repo.MakeBooking)
		mux.Post("/make-booking", controllers.Repo.MakeBookingPost)
		mux.Get("/confirm-booking", controllers.Repo.ConfirmBooking)
		mux.Post("/confirm-booking", controllers.Repo.ConfirmBookingPost)
	})
	mux.Get("/logout", controllers.Repo.Logout)

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

	fileServer := http.FileServer(http.Dir("./resources/"))
	mux.Handle("/resources/*", http.StripPrefix("/resources", fileServer))
	return mux
}
