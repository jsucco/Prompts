package controllers

import (
	"net/http"
	"viewmodels"
	"text/template"
)

type gmController struct {
	template *template.Template
}

func (this *gmController) get(w http.ResponseWriter, req *http.Request) {
	vm := viewmodels.GetGM()

	w.Header().Add("Content Type", "text/html")

	this.template.Execute(w, vm)
}