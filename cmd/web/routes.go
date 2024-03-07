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
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet,"/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet,"/vlog/view/:id",dynamic.ThenFunc(app.vlogView))
	router.Handler(http.MethodGet,"/vlog/create",dynamic.ThenFunc(app.vlogCreate))
	router.Handler(http.MethodPost,"/vlog/create", dynamic.ThenFunc(app.vlogCreatePost))
	router.Handler(http.MethodPut,"/vlog/update",dynamic.ThenFunc(app.vlogUpdate))
	router.Handler(http.MethodDelete,"/vlog/delete",dynamic.ThenFunc(app.vlogDelete))
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	// Return the 'standard' middleware chain followed by the servemux.
	return standard.Then(router)
}
