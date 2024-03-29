package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

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

	fieldErrors := make(map[string]string)

	// Check that the title value is not blank and is not more than 100
	// characters long. If it fails either of those checks, add a message to the
	// errors map using the field name as the key.
	if strings.TrimSpace(title) == "" {
		fieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		fieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	// Check that the description value isn't blank.
	if strings.TrimSpace(description) == "" {
		fieldErrors["description"] = "This field cannot be blank"
	}

	// Check if a file is uploaded
	if header.Filename == "" {
		fieldErrors["file"] = "No file uploaded"
	} else {
		// Create a unique file name for the uploaded file
		// For simplicity, just using the original filename here. Consider generating a unique name.
		// Ensure your application saves the file in an appropriate directory with proper permissions
		photoFileName := header.Filename
		filePath := "./ui/static/img/" + photoFileName

		// Check file size
		if header.Size == 0 {
			fieldErrors["file"] = "Uploaded file is empty"
		}

		if len(fieldErrors) == 0 {
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
		}
	}

	if len(fieldErrors) > 0 {
		// Convert fieldErrors map to JSON for better formatting
		fieldErrorsJSON, err := json.Marshal(fieldErrors)
		if err != nil {
			app.serverError(w, err)
			return
		}
		// Set appropriate content type for response
		w.Header().Set("Content-Type", "application/json")
		// Write JSON response with field errors
		w.WriteHeader(http.StatusBadRequest) // Or appropriate HTTP status code
		w.Write(fieldErrorsJSON)
		return
	}

	// Insert the vlog details into the database
	id, err := app.vlogs.Insert(user_id, title, description, header.Filename, 0, 0)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "vlog created successfuly")

	http.Redirect(w, r, fmt.Sprintf("/vlog/view/%d", id), http.StatusSeeOther)
}
func (app *application) vlogUpdate(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	idStr := params.ByName("vlog_id") // Change "id" to "vlog_id" to match the route definition
	// Assuming your route parameter is named "id"
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	vlog, err := app.vlogs.Get(id) // Assuming app.vlogs.GetByID fetches the vlog details by ID
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Vlog = vlog // Populate the vlog details
	app.render(w, http.StatusOK, "update.html", data)
}

func (app *application) vlogUpdatePost(w http.ResponseWriter, r *http.Request) {
	// Get the vlog ID from the request parameters
	params := httprouter.ParamsFromContext(r.Context())
	idStr := params.ByName("vlog_id") // Make sure this matches your route parameter name
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	// Fetch existing vlog details in case we need to fallback to original values
	existingVlog, err := app.vlogs.Get(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Parse form data with a 10 MB limit
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Retrieve the updated title and description from the form, defaulting to existing values if not provided
	newTitle := r.FormValue("title")
	if newTitle == "" {
		newTitle = existingVlog.Title
	}

	newDescription := r.FormValue("description")
	if newDescription == "" {
		newDescription = existingVlog.Description
	}

	// Initialize photoFileName with the existing photo file name
	photoFileName := existingVlog.PhotoFile

	// Process file upload if a file is provided
	file, header, err := r.FormFile("photoFile")
	if err == nil && header != nil { // Only proceed if a file was actually uploaded
		defer file.Close()
		photoFileName = header.Filename
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
	} // If an error occurred because no file was uploaded, we simply retain the existing photo file

	// Update the vlog in the database
	err = app.vlogs.Update(id, newTitle, newDescription, photoFileName)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect to the vlog view page
	http.Redirect(w, r, fmt.Sprintf("/vlog/view/%d", id), http.StatusSeeOther)
}

func (app *application) vlogDelete(w http.ResponseWriter, r *http.Request) {
    // Extract the vlog ID from the URL parameters
    params := httprouter.ParamsFromContext(r.Context())
    vlogIDStr := params.ByName("vlog_id")

    // Validate the vlog ID
    vlogID, err := strconv.Atoi(vlogIDStr)
    if err != nil || vlogID < 1 {
        app.notFound(w)
        return
    }

    // Delete the vlog from the database
    err = app.vlogs.Delete(vlogID)
    if err != nil {
        app.serverError(w, err)
        return
    }
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

