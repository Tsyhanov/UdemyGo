package handlers

import (
	"net/http"
	"udemygo/basicwebapp/pkg/render"
)

func Home(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "This is home page")
	render.RenderTemplate(w, "home.page.tmpl")
}

func About(w http.ResponseWriter, r *http.Request) {
	//sum := AddValue(2, 2)
	//_, _ = fmt.Fprintf(w, fmt.Sprintf("This is a AddValue result: %d", sum))
	render.RenderTemplate(w, "about.page.tmpl")
}
