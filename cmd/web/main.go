package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"h26.io/platform/pkg/config"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	templateCache map[string]*template.Template
}

func main() {
	err := realMain()
	if err != nil {
		os.Exit(2)
	}
}

func realMain() error {
	// Check and create structure with configuration
	conf.InitConfig()

	// default port
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// Use log.New() to create a logger for writing information messages.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	// initialize database connection pool
	db, err := openDB(genDSN())
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Initialize a new template cache...
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// Initialize a new instance of application containing the dependencies.
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Write messages using the two new loggers, instead of the standard logger.
	infoLog.Printf("Starting server on %s", *addr)

	err = srv.ListenAndServe()
	errorLog.Fatal(err)

	return nil
}

func openDB(dsn string) (*sql.DB, error) {
	// this is only initialization, we are not calling DB yet.
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	// it's a good idea to check if we can ping DB
	//err = db.Ping()
	//if err != nil {
	//		return nil, err
	//	}
	return db, nil
}

func genDSN() string {
	parseTime := true
	return fmt.Sprintf(
		"%v:%v@/%v?parseTime=%v",
		conf.Config.DbUsername,
		conf.Config.DbPassword,
		conf.Config.DbName,
		parseTime,
	)
}
