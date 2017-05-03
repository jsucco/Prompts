package controllers

import (
	"net/http"
	"text/template"
	"github.com/gorilla/mux"
)

var tpl *template.Template

func Register(templates *template.Template)  {
	tpl = templates
	router := mux.NewRouter()

	hc := new(homeController)
	hc.template = templates.Lookup("home.gohtml")
	router.HandleFunc("/home", hc.get)

	gm := new(gmController)
	gm.template = templates.Lookup("gm.gohtml")
	router.HandleFunc("/gm", gm.handle)

	router.HandleFunc("/", idx)

	http.Handle("/", router)

	http.Handle("/public/", http.FileServer(http.Dir(".")))

}

func idx(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w,"home.gohtml", nil)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

}
