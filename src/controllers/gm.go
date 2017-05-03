package controllers

import (
	"net/http"
	"viewmodels"
	"text/template"
	"strconv"
)

type gmController struct {
	template *template.Template
}

func (this *gmController) handle(w http.ResponseWriter, req *http.Request) {
	vm := viewmodels.GetGM()

	if req.Method == "POST" {
		var pagertext = req.FormValue("PagerText")
		j, err := strconv.Atoi(req.FormValue("SelectedPrompt"))
		if err == nil {
			if pagertext == "next" {
				vm.SelectedPrompt = j + 1
			} else if pagertext == "prev" {
				vm.SelectedPrompt = j - 1
			}
		}
	}

	w.Header().Add("Content Type", "text/html")

	this.template.Execute(w, vm)
}

func (this *gmController) post(w http.ResponseWriter, req *http.Request) {

}