package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"udemygo/basicwebapp/pkg/config"
	"udemygo/basicwebapp/pkg/models"
)

var app *config.AppConfig

//NewTemplate sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {

	var tc map[string]*template.Template
	//get the template cache from app config
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	//get requested template
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td)

	_ = t.Execute(buf, td)

	//render the template
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	//myCache := make(map[string]*template.Template)
	myCache := map[string]*template.Template{}

	//get all files *.page.tmpl
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		log.Println("error reading templates")
		return myCache, err
	}
	//range through all files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			log.Println("error parse template files")
			return myCache, err
		}

		//check for the layouts
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			log.Println("error parse layout files")
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				log.Println("error parse layout files")
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}
