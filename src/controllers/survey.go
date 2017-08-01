package controllers

import (
	"net/http"
	"text/template"
	"viewmodels"
	_"strings"
	_"controllers/util"
	"errors"
	"models"
	"time"
)

type surveyController struct {
	template *template.Template
}

var (
	vm viewmodels.SurveyModel
    sm models.Survey
	err error
)

func (this *surveyController) handle(w http.ResponseWriter, req *http.Request) {

	if Authenticated(req) == false {
		http.Redirect(w, req, "/login?returnURL=" + req.RequestURI, 302)
		return
	}
	vm = viewmodels.GetEmptyGMSurveyModel()

	if req.Method != "POST" {
		sm, err = getSurveyCtx(req)
		if err != nil {
			vm.ErrorMessage = err.Error()
			sm = models.AssembleGM()
		}
	} else {
		sm, _ = models.DeserializeBuffer(req.FormValue("ModelBuffer"))

		error := handleUserResponses(req)
		if error != nil {
			vm.ErrorMessage = error.Error()
		}
		//validateAnswers(req)
	}
	vm.SetPrompts(req, sm)
	vm.SetModelBuffer(sm)
	w.Header().Add("Content Type", "text/html")

	this.template.Execute(w, vm)
}

func getSurveyCtx(req *http.Request) (models.Survey, error) {
	q_params := req.URL.Query()
	if q_params != nil {
		surveyid := q_params.Get("surveyid")
		if len(surveyid) > 0 {
			var sessionid = models.ReadSessionCookie(req)
			if len(sessionid) > 0 {
				sess, err := models.GetUserSession(sessionid, req)
				if err == nil {
					s := models.Survey{}
					err := s.LoadSurvey(sess.OrganizationKey, surveyid, req)
					if err == nil {
						return s, nil
					} else {
						return models.Survey{}, errors.New(err.Error() + ": " + surveyid + ": " + sess.OrganizationKey)
					}
				} else {
					return models.Survey{}, err
				}
			} else {
				return models.Survey{}, errors.New("sessioniid cannot be length 0.")
			}

		} else {
			return models.Survey{}, errors.New("failed to get query string var surveyid.")
		}
	}
	return models.AssembleGM(), nil
}

func handleUserResponses(req *http.Request) error {
	var sess models.Session
	var err error
	if len(sm.Prompts) > 0 {
		sm.MapAllResponses(req)
		sm.Completed = true
		sm.Finished = time.Now().Local()
		var sessionid = models.ReadSessionCookie(req)
		if len(sessionid) > 0 {
			sess, err = models.GetUserSession(sessionid, req)
			if err == nil {
				if err = sm.SaveSurvey(sess.OrganizationKey, req); err == nil {
					if err = sm.AddNewAsset(sess.OrganizationKey, req); err == nil {
						return nil
					}
				}
			}
		}
	} else {
		err = errors.New("at least one prompt is required.")
	}
	return err
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

func validateSingleAnswer(q models.Question) (bool, string) {
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
