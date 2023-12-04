package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
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
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/about-us", app.aboutUs)
	mux.HandleFunc("/damzxyno", app.damzxyno)
	mux.HandleFunc("/panic", app.panic)
	mux.Handle("/static/", http.StripPrefix("/static", fS))
	mux.HandleFunc("/snippet/view", app.view)
	mux.HandleFunc("/err", func(w http.ResponseWriter, r *http.Request) {
		app.serverError(w, (errors.New("This is an information")))
	})
	mux.HandleFunc("/snippet/create", func(a http.ResponseWriter, b *http.Request) {
		if b.Method != http.MethodPost {
			a.Header().Set("Allow", http.MethodPost)
			a.WriteHeader(http.StatusMethodNotAllowed)
			a.Write([]byte("Aurevoir!"))
			return
		}
		a.Write([]byte("La loopy"))
	})
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

func (g application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.infoLog.Println("This is an information")
	//infoLog.Fatal("This is an error")
	g.errLog.Println("This ia an information II")
	g.errLog.Fatal("This is an error II")
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
