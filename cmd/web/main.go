package main

import (

	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"

	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"twitbox.vedantkugaonkar.net/internal/model"
)

type application struct {
	infoLog        *log.Logger
	errorLog       *log.Logger
	twits          *model.TwitModel
	tempCache      map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
	users          *model.UserModel
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	// dsn=data source name,contains all the info req for establishing a connection
	//to the database.Just a string with all info for the driver
	dsn := flag.String("dsn", os.Getenv("TWITBOX_DB_DSN"), "MySQL data source name")
	flag.Parse()

	if *dsn == "" {
		// fallback for local development
		*dsn = "twitbox:rama2804@tcp(localhost:3306)/twitbox?parseTime=true&multiStatements=true"

	}

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	//imp to defer the db before main finishes
	defer db.Close()
	tempCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}
	formDecoder := form.NewDecoder()
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true
        sessionManager.Cookie.HttpOnly = true
        sessionManager.Cookie.SameSite = http.SameSiteLaxMode
        sessionManager.Cookie.Path = "/"
        sessionManager.Cookie.Name = "session"
        //sessionManager.Cookie.Domain = "twitbox.app"

	

	ap := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		twits:          &model.TwitModel{DB: db},
		tempCache:      tempCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
		users:          &model.UserModel{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		//calling the ap.routes() to get the servemux containing our routes
		Handler:      ap.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting a server on %s", *addr)
	err = srv.ListenAndServe()

	errorLog.Fatal(err)

}
