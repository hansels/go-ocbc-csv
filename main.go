package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/csv", func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(32 << 20)

		file, handler, err := r.FormFile("file")
		if err != nil {
			fmt.Println("Error Retrieving the File")
			fmt.Println(err)
			return
		}
		defer file.Close()

		fmt.Printf("Uploaded File: %+v\n", handler.Filename)
		fmt.Printf("File Size: %+v\n", handler.Size)
		fmt.Printf("MIME Header: %+v\n", handler.Header)

		tempFile, err := ioutil.TempFile("temp-images", "upload-*.csv")
		if err != nil {
			fmt.Println(err)
		}
		defer tempFile.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}

		tempFile.Write(fileBytes)

		fmt.Fprintf(w, "Successfully Uploaded File\n")
	})

	r.Get("/csv", func(w http.ResponseWriter, r *http.Request) {
		tempFile, err := ioutil.ReadFile("temp-images/upload-*.csv")
		if err != nil {
			fmt.Println(err)
		}

		fmt.Fprint(w, string(tempFile))
	})

	http.ListenAndServe(":3000", r)
}
