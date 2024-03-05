package main

import (
	"net/http"

	"github.com/justinas/alice"
)

// the route() method returns servermux containing our application routes
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui//static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/vlog/view", app.vlogView)
	mux.HandleFunc("/vlog/create", app.vlogCreate)
	mux.HandleFunc("/vlog/update", app.vlogUpdate)
	mux.HandleFunc("/vlog/delete", app.vlogDelete)
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	// Return the 'standard' middleware chain followed by the servemux.
	return standard.Then(mux)
}
