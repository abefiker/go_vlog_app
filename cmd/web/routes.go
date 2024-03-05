package main

import "net/http"

//the route() method returns servermux containing our application routes
func (app *application) routes() http.Handler{
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui//static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/vlog/view", app.vlogView)
	mux.HandleFunc("/vlog/create", app.vlogCreate)
	mux.HandleFunc("/vlog/update", app.vlogUpdate)
	mux.HandleFunc("/vlog/delete", app.vlogDelete)
	return secureHeaders(mux)
}
