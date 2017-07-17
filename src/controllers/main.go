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
	router.HandleFunc("/home", hc.get)
	router.HandleFunc("/login", hc.login)

	sy := new(surveyController)
	sy.template = templates.Lookup("survey.gohtml")
	router.HandleFunc("/survey", sy.handle)

	router.HandleFunc("/", idx)

	http.Handle("/", router)

	http.Handle("/public/", http.FileServer(http.Dir(".")))

}

func idx(w http.ResponseWriter, r *http.Request) {

	if Authenticated(r) == false {
		http.Redirect(w, r, "/login?returnURL=" + r.RequestURI, 302)
		return
	}

	err := tpl.ExecuteTemplate(w,"home.gohtml", nil)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

}
