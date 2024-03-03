package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// the folowing application struct purpose is for dependency injection
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {

	addr := flag.String("addr", ":8080", "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui//static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/vlog/view", app.vlogView)
	mux.HandleFunc("/vlog/create", app.vlogCreate)
	mux.HandleFunc("/vlog/update", app.vlogUpdate)
	mux.HandleFunc("/vlog/delete", app.vlogDelete)

	srv := http.Server{
		Addr:     *addr,
		Handler:  mux,
		ErrorLog: errorLog,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}