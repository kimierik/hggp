package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"compress/gzip"
	"strings"
)

var GLOBAL_TEMPLATES *template.Template 

func reset_templates(w http.ResponseWriter, r *http.Request){
    tmpl := template.Must(template.ParseGlob("templates/*.html"))

	GLOBAL_TEMPLATES = tmpl
}

func addHandles(){
	http.HandleFunc("/", HandleIndex)
	http.HandleFunc("/fresh", reset_templates)
}

type gzipResponseWriter struct {
    http.ResponseWriter
    Writer *gzip.Writer
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
    return w.Writer.Write(b)
}


// gzipHandler wraps a handler and gzips the response if the client supports it
func gzipHandler(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Only gzip if client supports it
        if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
            h.ServeHTTP(w, r)
            return
        }

        w.Header().Set("Content-Encoding", "gzip")
        gz := gzip.NewWriter(w)
        defer gz.Close()

        gzw := &gzipResponseWriter{ResponseWriter: w, Writer: gz}
        h.ServeHTTP(gzw, r)
    })
}


func main() {
	fmt.Println("frontend start")
    tmpl := template.Must(template.ParseGlob("templates/*.html"))

	GLOBAL_TEMPLATES = tmpl

	fs := http.FileServer(http.Dir("./static"))
    http.Handle("/static/", gzipHandler(http.StripPrefix("/static/", fs)))

	addHandles()

    log.Fatal(http.ListenAndServe(":8080", nil))
}
