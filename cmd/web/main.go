package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"snippetbox.damzxyno.net/internal/models"
)

type Configuration struct {
	address string
}

type application struct {
	infoLog       *log.Logger
	errLog        *log.Logger
	snippets      *models.SnippetModel
	htmlTemplates map[string]*template.Template
}

func (app *application) routes() http.Handler {
	fS := http.FileServer(http.Dir("./ui/static/"))
	mux := httprouter.New()
	mux.HandlerFunc(http.MethodGet, "/", app.home)
	mux.HandlerFunc(http.MethodGet, "/about-us", app.aboutUs)
	mux.HandlerFunc(http.MethodGet, "/damzxyno", app.damzxyno)
	mux.HandlerFunc(http.MethodGet, "/panic", app.panic)
	mux.Handler(http.MethodGet, "/static/", http.StripPrefix("/static", fS))
	mux.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.viewById)
	mux.HandlerFunc(http.MethodGet, "/snippet/create", app.createSnippet)
	mux.HandlerFunc(http.MethodPost, "/snippet/create", app.createSnippetPost)
	standard := alice.New(app.handlePanic, app.logRequest, secureHeader)
	return standard.Then(mux)
}

func resolveAllTemplate() (map[string]*template.Template, error) {
	cache := make(map[string]*template.Template)

	if files, err := filepath.Glob("./ui/html/pages/*.tmpl"); err != nil {
		return nil, err
	} else {
		for _, file := range files {
			baseTemp, baseTempErr := template.ParseFiles("./ui/html/base.tmpl")
			if baseTempErr != nil {
				return nil, baseTempErr
			}
			baseAndPartialTemp, baseAndPartialTempErr := baseTemp.ParseGlob("./ui/html/partials/*.tmpl")
			if baseAndPartialTempErr != nil {
				return nil, baseTempErr
			}
			if tem, err := baseAndPartialTemp.ParseFiles(file); err != err {
				return nil, err
			} else {
				cache[filepath.Base(file)] = tem
			}
		}
	}
	return cache, nil
}
func main() {
	var config Configuration
	var f, err0 = os.OpenFile("./tmp/info.log", os.O_RDWR|os.O_CREATE, 0666)
	if err0 != nil {
		log.Fatal(err0)
	}
	defer f.Close()

	app := &application{
		infoLog: log.New(f, "[INFO]\t", log.Ldate|log.Ltime),
		errLog:  log.New(f, "[ERR]\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
	templates, templateErr := resolveAllTemplate()
	if templateErr != nil {
		errorsString := fmt.Sprintf("%v, %v", templateErr.Error(), debug.Stack())
		app.errLog.Output(2, errorsString)
	}
	app.htmlTemplates = templates

	db, dbErr := openDb("web:password@/snippetbox?parseTime=true")
	if dbErr != nil {
		errorNote := fmt.Sprintf("%s%n %s", dbErr.Error(), debug.Stack())
		app.errLog.Output(2, errorNote)
	}

	app.snippets = &models.SnippetModel{DB: db}

	defer db.Close()
	flag.StringVar(&config.address, "addr", "localhost:4000", "HTTP Netword Address")
	flag.Parse()

	log.Print("Http Server Starting. . . .")
	server := &http.Server{
		Addr:     config.address,
		ErrorLog: app.errLog,
		Handler:  app.routes(),
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func openDb(dns string) (*sql.DB, error) {
	db, dbErr := sql.Open("mysql", dns)
	if dbErr != nil {
		return nil, dbErr
	}
	if isPingableError := db.Ping(); isPingableError != nil {
		return nil, isPingableError
	}
	return db, nil
}
