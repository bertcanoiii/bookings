package main

import (
	"github.com/bertcanoiii/bookings/pkg/config"
	"github.com/bertcanoiii/bookings/pkg/handlers"
	"github.com/bertcanoiii/bookings/pkg/render"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"time"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

// main is the main application function
func main() {

	//Here we create a variable of type appconfig
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	//We create our TemplateCache and assign it to tc
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	//Here we assign tc to app.TemplateCache
	//This will make it so that the app only creates the TemplateCache once
	//Instead of on every request
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	//This sets app in render.go to the app.config templateCache
	render.NewTemplates(&app)

	fmt.Println(fmt.Sprintf("Starting application on port: %s", portNumber))

	//This creates a new server
	//In it we send the app config to the routes.go > routes() func we made
	srv := &http.Server{
		Addr: portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}