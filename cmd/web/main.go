package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/SmoothWay/snippetBox/pkg/models/postgresql"
	_ "github.com/lib/pq"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *postgresql.SnippetModel
	templateCache map[string]*template.Template
}

const (
	host     = "localhost"
	port     = 5432
	user     = "weblocalhost"
	password = "insaneEra99"
	dbname   = "snippetbox"
)

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")
	// dsn := flag.String("dsn", "weblocalhost:patrickBateman99@http://localhost:5432/snippetbox?parseTime=true", "postgresql connection")
	flag.Parse()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(psqlInfo)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	tempalteCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &postgresql.SnippetModel{DB: db},
		templateCache: tempalteCache,
	}

	infoLog.Printf("Starting server on %s\n", *addr)

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	err = srv.ListenAndServe()
	errorLog.Fatal(err)

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
