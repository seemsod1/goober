package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/seemsod1/goober/controllers"
	"github.com/seemsod1/goober/helpers/render"
	"github.com/seemsod1/goober/initializers"
	models "github.com/seemsod1/goober/models/app_models"
	"github.com/seemsod1/goober/models/entities"
	"log"
	"net/http"
	"os"
	"time"
)

func init() {
	initializers.LoadConfig()
}

var app models.AppConfig
var session *scs.SessionManager

func main() {

	app.DB = initializers.ConnectToDatabase()
	initializers.SyncDB(app.DB)
	//initializers.Migration(app.DB)
	//Production
	gob.Register(entities.User{})
	gob.Register(entities.RentInfo{})
	gob.Register(entities.CarHistory{})
	gob.Register(entities.UserHistory{})
	gob.Register(uuid.UUID{})

	app.MailChan = make(chan models.MailData)

	defer close(app.MailChan)
	listenForMail()

	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}
	app.TemplateCache = tc
	app.UseCache = false

	app.Scheduler, err = gocron.NewScheduler()
	if err != nil {
		log.Fatal("cannot create scheduler")
	}
	app.Scheduler.Start()
	repo := controllers.NewRepo(&app)
	controllers.NewControllers(repo)
	render.NewRenderer(&app)

	portNumber := ":" + os.Getenv("PORT")
	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
