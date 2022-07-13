package main

import (
	"log"
	"net/http"
	"udemygo/basicwebapp/pkg/config"
	"udemygo/basicwebapp/pkg/handlers"
	"udemygo/basicwebapp/pkg/render"
)

func main() {

	var app config.AppConfig

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}
	app.TemplateCache = tc
	app.UseCache = false
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	//give access to app config from the render package
	render.NewTemplates(&app)

	http.HandleFunc("/", handlers.Repo.Home)
	http.HandleFunc("/about", handlers.Repo.About)

	_ = http.ListenAndServe("127.0.0.1:8080", nil)
}
