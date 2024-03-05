package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// the route() method returns servermux containing our application routes
func (app *application) routes() http.Handler {
	// mux := http.NewServeMux()
	router := httprouter.New()
	fileServer := http.FileServer(http.Dir("./ui//static/"))
	router.Handler(http.MethodGet,"/static/", http.StripPrefix("/static", fileServer))

	router.HandlerFunc(http.MethodGet,"/", app.home)
	router.HandlerFunc(http.MethodGet,"/vlog/view/:id", app.vlogView)
	router.HandlerFunc(http.MethodGet,"/vlog/create", app.vlogCreate)
	router.HandlerFunc(http.MethodPost,"/vlog/create", app.vlogCreatePost)
	router.HandlerFunc(http.MethodPut,"/vlog/update", app.vlogUpdate)
	router.HandlerFunc(http.MethodDelete,"/vlog/delete", app.vlogDelete)
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	// Return the 'standard' middleware chain followed by the servemux.
	return standard.Then(router)
}
