package render

import (
	"github.com/bertcanoiii/bookings/pkg/config"
	"github.com/bertcanoiii/bookings/pkg/models"
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

//This will hold any functions we want to use in the tmpl files
var functions = template.FuncMap{

}

//This will hold our app.config templatecache
var app *config.AppConfig

// NewTemplates sets app to templatecache
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// RenderTemplate renders templates using html/template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {

	//Create templateCache variable
	var templateCache map[string]*template.Template

	//UseCache is set to true in the app.Config file
	//We create the template cache once on startup
	//Otherwise (if in development mode) we recreate the cache every request
	if app.UseCache {
		templateCache = app.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()
	}

	//pageTemplate = templateCache with index tmpl
	pageTemplate, pageExists := templateCache[tmpl]
	if !pageExists {
		log.Fatal("Couldn't get template cache")
	}

	//create byteBuffer variable
	byteBuffer := new(bytes.Buffer)

	td = AddDefaultData(td)

	//Execute pageTemplate and store into byteBuffer
	_ = pageTemplate.Execute(byteBuffer, td)

	//Write byteBuffer contents to response writer
	_, err := byteBuffer.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
	}

}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {

	//Creates a map which will hold all our templates
	myCache := map[string]*template.Template{}

	//Finds everything in the templates folder that ends in .page.tmpl
	//Assigns that to "pages"
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	//For each page in pages
	//name = actual filepath name
	//create a template set, ts
	//ts holds the template name, any functions we want to use, and parsed page data
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		//If any files in the template folder match ".layout.tmpl" find them
		//Add to variable "matches"
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		//If there's anything in matches
		//Add to template set
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		//Add everything we just did to myCache map
		myCache[name] = ts

	}

	return myCache, nil

}