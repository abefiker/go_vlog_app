package main

import (
	"net/http"
	"text/template"
)

func (app *application)home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/home.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}
func (app *application)vlogView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("vlogs will be displayed here"))
}
func (app *application)vlogCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("form to create vlog will displayed here"))
}
func (app *application)vlogUpdate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("update functionality"))
}
func (app *application)vlogDelete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("delete functionality"))
}
