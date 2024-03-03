package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)
func uploadFile(w http.ResponseWriter, r *http.Request, fileType string) {
    // Check if the request method is POST
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Retrieve the file from the request
    file, header, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "Error retrieving file", http.StatusInternalServerError)
        return
    }
    defer file.Close()

    // Determine the upload directory based on the file type
    var uploadDir string
    if fileType == "video" {
        uploadDir = "./uploads/videos/"
    } else if fileType == "image" {
        uploadDir = "./uploads/images/"
    } else {
        http.Error(w, "Invalid file type", http.StatusBadRequest)
        return
    }

    // Create the upload directory if it doesn't exist
    if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
        os.MkdirAll(uploadDir, 0755)
    }

    // Create a unique file name
    filename := uploadDir + header.Filename

    // Create the file
    outputFile, err := os.Create(filename)
    if err != nil {
        http.Error(w, "Error saving file", http.StatusInternalServerError)
        return
    }
    defer outputFile.Close()

    // Copy the uploaded file to the newly created file
    _, err = io.Copy(outputFile, file)
    if err != nil {
        http.Error(w, "Error saving file", http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "File uploaded successfully")
}
