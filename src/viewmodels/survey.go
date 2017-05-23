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
	ModelBuffer string
	SurveyName string
	PrevPrompt int
}

func GetEmptyGMSurveyModel() SurveyModel {

	model := SurveyModel{
		Title: "General Manager",
		SelectedPrompt: 0,
		PrevPrompt: 0,
		ModelBuffer: "",
		SurveyName: "",
	}

	return model
}

func (survey *SurveyModel) SetQuestions(req *http.Request, model Surveys.Survey)  {

	p := model.Prompts[survey.SelectedPrompt]

	survey.LastPrompt = 4
	survey.Questions = p.Questions

}

func (survey *SurveyModel) SetModelBuffer(model Surveys.Survey) {
	survey.ModelBuffer = model.ToBase64()
}
