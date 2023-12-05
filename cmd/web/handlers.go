package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"snippetbox.damzxyno.net/internal/models"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	tmp := app.htmlTemplates["home.tmpl"]
	err2 := tmp.ExecuteTemplate(w, "base", nil)
	if err2 != nil {
		app.serverError(w, err2)
	}
}

func (app *application) aboutUs(w http.ResponseWriter, r *http.Request) {
	tmp := app.htmlTemplates["about-us.tmpl"]
	tmp.ExecuteTemplate(w, "base", nil)
}

func (app *application) damzxyno(w http.ResponseWriter, r *http.Request) {
	log.Print(fmt.Sprintf("SINGLE: %s", w.Header().Get("Weight")))
	log.Print(fmt.Sprintf("Multiple: %v", w.Header().Values("Weight")))
	w.Write([]byte("{\n\t\"Name\":\"Damilola\", \n\t\"Occupation\":\"Computer Scientist\"\n}"))
}

func (app *application) panic(w http.ResponseWriter, r *http.Request) {
	panic("This is an intentional panic!")
}

type viewSnippets struct {
	Snippets []*models.Snippet
}

func (app *application) viewById(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	val, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	snip, _ := app.snippets.Get(val)
	tmp := app.htmlTemplates["view.tmpl"]
	tmp.ExecuteTemplate(w, "base", snip)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	tmp := app.htmlTemplates["create.tmpl"]
	tmp.ExecuteTemplate(w, "base", nil)
}

func (app *application) createSnippetPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
		return
	}
	title := r.PostForm.Get("title")
	content := r.PostFormValue("content")
	expires, err2 := strconv.Atoi(r.PostForm.Get("expires"))
	if err2 != nil {
		app.serverError(w, err2)
		return
	}
	id, err3 := app.snippets.Insert(title, content, expires)
	if err3 != nil {
		app.serverError(w, err3)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
