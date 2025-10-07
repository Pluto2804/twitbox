package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"twitbox.vedantkugaonkar.net/internal/model"
)

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	twits    *model.TwitModel
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

	ap := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		twits:    &model.TwitModel{DB: db},
	}

	srv := http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		//calling the ap.routeMux() to get the servemux containing our routes
		Handler: ap.routeMux(),
	}

	infoLog.Printf("Starting a server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)

}
