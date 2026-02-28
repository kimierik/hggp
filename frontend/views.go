package main

import (
	"net/http"
)


func HandleIndex(w http.ResponseWriter, r *http.Request) {
	//hanle /

	// Execute the template named "index.html"
	d:= struct {
		PageTitle string 
	}{
		PageTitle: "mainpage",
	}


	if err := GLOBAL_TEMPLATES.ExecuteTemplate(w, "index.html", d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandleProfile(w http.ResponseWriter, r *http.Request) { }
