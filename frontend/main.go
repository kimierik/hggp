package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
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

func main() {
	fmt.Println("frontend start")
    tmpl := template.Must(template.ParseGlob("templates/*.html"))

	GLOBAL_TEMPLATES = tmpl

	fs := http.FileServer(http.Dir("./static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

	addHandles()

    log.Fatal(http.ListenAndServe(":8080", nil))
}
