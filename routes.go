package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"help/controllers"
	"net/http"
)

func routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", controllers.HomePage)
	mux.Get("/about", controllers.About)
	mux.Get("/cars", controllers.CarsPage)
	mux.Get("/login", controllers.LoginPage)

	fileServer := http.FileServer(http.Dir("./resources/"))
	mux.Handle("/resources/*", http.StripPrefix("/resources", fileServer))
	return mux
}
