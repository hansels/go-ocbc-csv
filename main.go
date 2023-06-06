package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

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

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}

		if _, err := os.Stat("csv"); os.IsNotExist(err) {
			os.Mkdir("csv", 0755)
		}

		fo, err := os.Create("csv/" + handler.Filename)
		if err != nil {
			fmt.Println(err)
		}
		defer fo.Close()

		_, err = fo.Write(fileBytes)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Fprintf(w, "Successfully Uploaded File\n")
	})

	r.Get("/csv", func(w http.ResponseWriter, r *http.Request) {
		filename := r.URL.Query().Get("fileName")
		tempFile, err := ioutil.ReadFile("csv/" + filename + ".csv")
		if err != nil {
			fmt.Println(err)
		}

		fmt.Fprint(w, string(tempFile))
	})

	http.ListenAndServe(":3000", r)
}
