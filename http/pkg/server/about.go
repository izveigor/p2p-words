package server

import (
	"net/http"
	"text/template"
)

func About(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("").ParseFiles(PATH_TO_STATIC+"base.html", PATH_TO_STATIC+"about.html")
	if err != nil {
		InternalServerError(w, err)
	}

	if err = tmpl.ExecuteTemplate(w, "base", nil); err != nil {
		InternalServerError(w, err)
	}
}
