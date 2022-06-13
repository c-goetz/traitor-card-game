package main

import (
	"bytes"
	"embed"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
)

//go:embed static/htmx.min.js
var htmx []byte

//go:embed templates
var templates embed.FS

func main() {
	tsFS, err := fs.Sub(templates, "templates")
	if err != nil {
		log.Fatal(err)
	}
	ts := template.Must(template.ParseFS(tsFS, "*.html"))
	mux := http.NewServeMux()
	mux.HandleFunc("/static/htmx.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/javascript")
		reader := bytes.NewReader(htmx)
		io.Copy(w, reader)
	})
	mux.HandleFunc("/index.html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		ts.ExecuteTemplate(w, "index.html", nil)
	})
	err = http.ListenAndServe(":8000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
