package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"help/controllers"
	"help/initializers"
	models "help/models/app_models"
	"log"
	"net/http"
	"os"
	"time"
)

func init() {
	initializers.LoadConfig()
	initializers.ConnectToDatabase()
	initializers.SyncDB()
}

var app models.AppConfig
var session *scs.SessionManager

func main() {

	//Production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := initializers.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.UseCache = false
	app.TemplateCache = tc

	controllers.SetAppForTemplate(&app)

	portNumber := ":" + os.Getenv("PORT")
	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
