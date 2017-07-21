package controllers

import (
	"net/http"
	"text/template"
	"viewmodels"
	"controllers/helpers"
	_"strings"
	_"controllers/util"
	"errors"
	"models"
)

type surveyController struct {
	template *template.Template
}

var (
	vm viewmodels.SurveyModel
    sm Surveys.Survey
	err error
)

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

		error := handleUserResponses(req)
		if error != nil {
			vm.ErrorMessage = error.Error()
		}
		//validateAnswers(req)
		//sm = Surveys.AssembleGM()
	}
	vm.SetPrompts(req, sm)
	vm.SetModelBuffer(sm)
	w.Header().Add("Content Type", "text/html")

	this.template.Execute(w, vm)
}

func handleUserResponses(req *http.Request) error {
	if len(sm.Prompts) > 0 {
		sm.MapAllResponses(req)
		var sessionid = models.ReadSessionCookie(req)
		if len(sessionid) > 0 {
			sm.SaveSurvey(sessionid)
		}
	} else {
		return errors.New("at least one prompt is required.")
	}
	return nil
}

func validateAnswers(req *http.Request) {
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

func validateSingleAnswer(q Surveys.Question) (bool, string) {
	if q.Required == true {
		//var answer = strings.ToUpper(strings.Trim(q.UserResponse.Content, " "))
		//if q.UserResponse.IsDate == true {
		//	if answer == "MM/DD/YYYY" || len(answer) == 0 {
		//		return false, "A valid date response is required."
		//	}
		//} else if q.UserResponse.IsString == true {
		//	if len(answer) == 0 {
		//		return false, "This question requires a response."
		//	}
		//}
	}
	return true, ""
}
