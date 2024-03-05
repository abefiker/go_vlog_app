package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
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
	id, err := strconv.Atoi(params.ByName("id"))
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
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "create.html", data)

}

func (app *application) vlogCreatePost(w http.ResponseWriter, r *http.Request) {
	// Make sure you parse the form including the file upload
	err := r.ParseMultipartForm(10 << 20) // For example, 10 MB limit
	if err != nil {
		app.serverError(w, err)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	user_id := 1 // This could also be fetched from session or context if needed

	// Handle file upload
	file, header, err := r.FormFile("photoFile")
	if err != nil {
		app.serverError(w, err)
		return
	}
	defer file.Close()

	// Create a unique file name for the uploaded file
	// For simplicity, just using the original filename here. Consider generating a unique name.
	// Ensure your application saves the file in an appropriate directory with proper permissions
	photoFileName := header.Filename
	filePath := "./ui/static/img/" + photoFileName
	out, err := os.Create(filePath)
	if err != nil {
		app.serverError(w, err)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Insert the vlog details into the database
	id, err := app.vlogs.Insert(user_id, title, description, photoFileName, 0, 0)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/vlog/view/%d", id), http.StatusSeeOther)
}

func (app *application) vlogUpdate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("update functionality"))
}
func (app *application) vlogDelete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("delete functionality"))
}
