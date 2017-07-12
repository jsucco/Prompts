package viewmodels

import (
	"controllers/helpers"
	"net/http"
)

type SurveyModel struct {
	Title string
	SelectedPrompt int
	LastPrompt int
	Questions []Surveys.Question
	Prompts []Surveys.Prompt
	ModelBuffer string
	SurveyName string
	PrevPrompt int
}

func GetEmptyGMSurveyModel() SurveyModel {

	model := SurveyModel{
		Title: "General Manager",
		SelectedPrompt: 1,
		PrevPrompt: 0,
		ModelBuffer: "",
		SurveyName: "",
	}

	return model
}

func (survey *SurveyModel) SetPrompts(req *http.Request, model Surveys.Survey) {

	survey.Prompts = model.Prompts
	var cnt = 0

	for i, p := range survey.Prompts {
		var all_p_qs = p.ListAllQuestions()
		cnt = cnt + len(all_p_qs)
		if len(all_p_qs) > 0 {
			survey.Prompts[i] = insertFollows(p, all_p_qs)
		}

	}

	if len(survey.Prompts) > 0 {
		survey.LastPrompt = survey.Prompts[len(survey.Prompts) - 1].Id
	}
}

func insertFollows(prompt Surveys.Prompt, q_adds []Surveys.Question) Surveys.Prompt {
	prompt.Questions = q_adds
	return prompt
}

func (survey *SurveyModel) SetQuestions(req *http.Request, model Surveys.Survey)  {

	p := model.Prompts[survey.SelectedPrompt]

	survey.LastPrompt = 4
	survey.Questions = p.Questions

}

func (survey *SurveyModel) SetModelBuffer(model Surveys.Survey) {
	survey.ModelBuffer = model.ToBase64()
}
