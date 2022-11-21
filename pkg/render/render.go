package render

import (
	"bytes"
	"github.com/ppichugin/booking-for-breakfast/pkg/config"
	"github.com/ppichugin/booking-for-breakfast/pkg/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}
var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	// add some data
	return td
}

// RenderTemplate renders templates using html.Template
func RenderTemplate(w http.ResponseWriter, templateName string, td *models.TemplateData) {

	var tc map[string]*template.Template
	if app.UseCache {
		// get the template cache from AppConfig
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplate()
	}

	// get requested template from cache
	t, ok := tc[templateName]
	if !ok {
		log.Fatal("Could not get template from teh template cache")
	}

	buf := new(bytes.Buffer)
	td = AddDefaultData(td)
	err := t.Execute(buf, td)
	if err != nil {
		log.Println(err)
	}

	// render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

func CreateTemplate() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
