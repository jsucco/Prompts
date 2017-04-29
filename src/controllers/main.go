package controllers

import (
	"net/http"
	"text/template"
	"github.com/gorilla/mux"

)

func Register(templates *template.Template)  {

	router := mux.NewRouter()

	hc := new(homeController)
	hc.template = templates.Lookup("home.gohtml")
	router.HandleFunc("/home", hc.get)

	gm := new(gmController)
	gm.template = templates.Lookup("gm.gohtml")
	router.HandleFunc("/gm", gm.get)

	//http.Handle("/", router)

	http.Handle("/public/", http.FileServer(http.Dir(".")))
}
