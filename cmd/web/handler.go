package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/abefiker/go_vlog_app/internal/models"
	"github.com/julienschmidt/httprouter"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	vlogs, err := app.vlogs.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(r)
	data.Vlogs = vlogs
	// Pass the data to the render() helper as normal.
	app.render(w, http.StatusOK, "home.html", data)

}
func (app *application) vlogView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("vlog_id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	vlog, err := app.vlogs.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	data := app.newTemplateData(r)
	data.Vlog = vlog
	app.render(w, http.StatusOK, "view.html", data)

}
func (app *application) vlogCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display the form for creating a new snippet..."))
	}
	
func (app *application) vlogCreatePost(w http.ResponseWriter, r *http.Request) {
	

	// Assuming a form field "fileType" is sent specifying whether it's a video or image
	user_id := 1
	title := "I've played Mafia game"
	description := "It's interesting and funny game, beside logics and reasoning also some human behavior will be revealed. However, if you don't control yourself you would waste too much time, so that is the bad side of it."
	photoFile := "love.jpg"
	views := 0
	likes := 0
	id, err := app.vlogs.Insert(user_id, title, description, photoFile, views, likes)

	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/vlog/view?id=%d", id), http.StatusSeeOther)
	// Call the uploadFile function with the appropriate file type
}

func (app *application) vlogUpdate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("update functionality"))
}
func (app *application) vlogDelete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("delete functionality"))
}
