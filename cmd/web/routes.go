package main

import (
	"github.com/bertcanoiii/bookings/pkg/config"
	"github.com/bertcanoiii/bookings/pkg/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

// routes creates our routes
func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)

	//This allows us to use static files from the static folder
	//First tell go where the static files are
	fileServer := http.FileServer(http.Dir("./static/"))
	//Tell the router to handle them
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}