package viewmodels

import (
	"models"
	"net/http"
)

type SurveyModel struct {
	Title string
	SelectedPrompt int
	LastPrompt int
	Questions []models.Question
	Prompts []models.Prompt
	ModelBuffer string
	SurveyName string
	PrevPrompt int
	ErrorMessage string
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

func (survey *SurveyModel) SetPrompts(req *http.Request, model models.Survey) {

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

func insertFollows(prompt models.Prompt, q_adds []models.Question) models.Prompt {
	prompt.Questions = q_adds
	return prompt
}

func (survey *SurveyModel) SetQuestions(req *http.Request, model models.Survey)  {

	p := model.Prompts[survey.SelectedPrompt]

	survey.LastPrompt = 4
	survey.Questions = p.Questions

}

func (survey *SurveyModel) SetModelBuffer(model models.Survey) {
	survey.ModelBuffer = model.ToBase64()
}
