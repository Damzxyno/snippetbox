package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"snippetbox.damzxyno.net/internal/models"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Content not found!"))
		return
	}

	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}
	tmp, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err2 := tmp.ExecuteTemplate(w, "base", nil)
	if err2 != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *application) aboutUs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed!", http.StatusMethodNotAllowed)
		return
	}

	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/pages/about-us.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/partials/sidebar.tmpl",
	}

	tmp, err := template.ParseFiles(files...)
	if err != nil {
		app.errLog.Fatal(err)
	}

	tmp.ExecuteTemplate(w, "base", nil)
}

func (app *application) damzxyno(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("CONTent-Type", "application/json")
	w.Header().Add("Weight", "10kg")
	w.Header().Add("Weight", "100lbs")
	w.Header()["MoNoPoLy"] = []string{"Vacation", "Jamainca"}
	w.Header().Set("Monopoly", "Background")
	w.Header().Set("MonoPOLY", "Frontground")
	w.Header().Add("MOnopoly", "Middleground")
	log.Print(fmt.Sprintf("SINGLE: %s", w.Header().Get("Weight")))
	log.Print(fmt.Sprintf("Multiple: %v", w.Header().Values("Weight")))
	w.Write([]byte("{\n\t\"Name\":\"Damilola\", \n\t\"Occupation\":\"Computer Scientist\"\n}"))
}

func (app *application) jiossm(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Welcome to JiossM ray!"))
}

type viewSnippets struct {
	Snippets []*models.Snippet
}

func (app *application) view(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/snippet/view" {

		snip, _ := app.snippets.Latest()
		//files := []string{
		//	"./ui/html/base.tmpl",
		//	"./ui/html/partials/nav.tmpl",
		//	"./ui/html/pages/view.tmpl",
		//}
		//
		//fmt.Println("JACK MAL")
		//tmpl, _ := template.ParseFiles(files...)

		tmpl := app.htmlTemplates["view.tmpl"]
		snippetz := viewSnippets{Snippets: snip}
		tmpErr := tmpl.ExecuteTemplate(w, "base", snippetz)
		if tmpErr != nil {
			fmt.Fprintf(w, "ERROR: %v", tmpl)
		}
		return
	}
	val, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	snip, _ := app.snippets.Get(val)
	log.Print(snip)
	fmt.Fprintf(w, "%+v\n, %T", *snip, *snip)
}
