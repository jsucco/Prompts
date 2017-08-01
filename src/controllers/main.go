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
	hc.loginTemplate = templates.Lookup("login.gohtml")

	router.HandleFunc("/login", hc.login)

	sy := new(surveyController)
	sy.template = templates.Lookup("survey.gohtml")
	router.HandleFunc("/survey", sy.handle)

	router.HandleFunc("/", hc.get)
	router.HandleFunc("/search", hc.search)

	http.Handle("/", router)

	http.Handle("/public/", http.FileServer(http.Dir(".")))

}
