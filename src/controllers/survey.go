package controllers

import (
	"net/http"
	"text/template"
	"strconv"
	"viewmodels"
	"controllers/helpers"
	"strings"
)

type surveyController struct {
	template *template.Template
}

var vm viewmodels.SurveyModel
var sm Surveys.Survey

func (this *surveyController) handle(w http.ResponseWriter, req *http.Request) {

	if Authenticated(req) == false {
		http.Redirect(w, req, "/login?returnURL=" + req.RequestURI, 302)
		return
	}

	vm = viewmodels.GetEmptyGMSurveyModel()

	if req.Method != "POST" {
		sm = Surveys.AssembleGM()
	} else {
		sm, _ = Surveys.DeserializeBuffer(req.FormValue("ModelBuffer"))
	}

	//handlePaging(req)
	mapAnswers(req)
	validateAnswers(req)
	//vm.SetQuestions(req, sm)
	vm.SetPrompts(req, sm)
	//vm.SetModelBuffer(sm)
	w.Header().Add("Content Type", "text/html")

	this.template.Execute(w, vm)


}

func handlePaging(req *http.Request) {
	if req.Method == "POST" {
		var pagertext = req.FormValue("PagerText")
		j, err := strconv.Atoi(req.FormValue("SelectedPrompt"))
		if err == nil {
			vm.PrevPrompt = j
			if pagertext == "next" {
				vm.SelectedPrompt = j + 1
			} else if pagertext == "prev" {
				vm.SelectedPrompt = j - 1
			}
		}
	}
}

func mapAnswers(req *http.Request) {
	if req.Method == "POST" {
		if len(sm.Prompts) >= vm.PrevPrompt {
			var p = sm.Prompts[vm.PrevPrompt]

			for i, s := range p.Questions {

				var inputval = req.FormValue(s.DataId)

				sm.Prompts[vm.PrevPrompt].Questions[i].UserResponse.Content = inputval

			}
		}

	}
}

func validateAnswers(req *http.Request) {
	if req.Method == "POST" {
		if len(sm.Prompts) >= vm.PrevPrompt {
			pt := req.FormValue("PagerText")

			if pt == "next" {
				var p = sm.Prompts[vm.PrevPrompt]

				for i, s := range p.Questions {
					v, msg := validateSingleAnswer(s)
					if v {
						sm.Prompts[vm.PrevPrompt].Questions[i].ErrorFlag = false
						sm.Prompts[vm.PrevPrompt].Questions[i].ErrorMessage = ""
					} else {
						sm.Prompts[vm.PrevPrompt].Questions[i].ErrorFlag = true
						sm.Prompts[vm.PrevPrompt].Questions[i].ErrorMessage = msg
						vm.SelectedPrompt = vm.PrevPrompt
					}
				}

			}
		}
	}
}

func validateSingleAnswer(q Surveys.Question) (bool, string) {
	if q.Required == true {
		var answer = strings.ToUpper(strings.Trim(q.UserResponse.Content, " "))
		if q.UserResponse.IsDate == true {
			if answer == "MM/DD/YYYY" || len(answer) == 0 {
				return false, "A valid date response is required."
			}
		} else if q.UserResponse.IsString == true {
			if len(answer) == 0 {
				return false, "This question requires a response."
			}
		}
	}
	return true, ""
}
