package main

import (
	"crypto/tls"
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
	dsn := flag.String("dsn", "frido:rama@2804@/twitbox?parseTime=true", "MySQL data source name")
	flag.Parse()

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

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

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
		//calling the ap.routeMux() to get the servemux containing our routes
		Handler:      ap.routeMux(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting a server on %s", *addr)
	err = srv.ListenAndServeTLS("tls/cert.pem", "tls/key.pem")
	errorLog.Fatal(err)

}
