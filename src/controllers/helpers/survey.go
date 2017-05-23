package Surveys

import (
	"fmt"
	"encoding/base64"
	"encoding/gob"
	"bytes"
	"strings"
)

type Survey struct {
	Title string
	Prompts []Prompt
	PromptCount int
	SelectedPrompt int
	LastPrompt int
}

type Prompt struct {
	Id int
	Questions []Question
	Followups []Question
	Complete bool
	FollowupsKey map[int]string
	SelectedPrompt int
	Title string
}

type Question struct {
	Id int
	QuestionText string
	DataId string
	UserResponse Response
	Required bool
	ErrorMessage string
	ErrorFlag bool
}

type Response struct {
	Content string
	Options []option
	Complete bool
	IsString bool
	IsNumeric bool
	IsDate bool
	IsSelect bool
}

type option struct {
	Value string
	Text string
	Selected bool
}

func DeserializeBuffer(buffer string) (Survey, bool) {
	if len(strings.Trim(buffer, " ")) == 0 {
		return Survey{}, false
	}
	return FromBase64(buffer), true
}

func AssembleGM() Survey {

	newSurvey := Survey{
		Title: "General Manager",
		Prompts: []Prompt{
			Prompt{
				Id: 1,
				Complete: false,
				FollowupsKey: nil,
				Title: "Location",
				Questions: []Question{
					Question{
						Id: 1,
						QuestionText: "What is your Street Address?",
						DataId: "StreetAddress",
						Required: true,
						UserResponse: Response{IsString: true, Content: ""},
					},
					Question{
						Id: 2,
						QuestionText: "What is your ZIP code?",
						DataId: "ZIP",
						Required: true,
						UserResponse: Response{IsString: true, },
					},
					Question{
						Id: 3,
						QuestionText: "City",
						DataId: "City",
						Required: true,
						UserResponse: Response{IsString: true, },
					},
					Question{
						Id: 4,
						QuestionText: "State",
						DataId: "State",
						Required: true,
						UserResponse: Response{IsString: true, },
					},
				},
			},
			Prompt{
				Id: 2,
				Complete: false,
				FollowupsKey: nil,
				Questions: []Question{
					Question{
						Id: 1,
						QuestionText: "Asset Install Date",
						DataId: "InstallDate",
						Required: true,
						UserResponse: Response{IsDate: true, },
					},
				},
			},
			Prompt{
				Id: 3,
				Complete: false,
				FollowupsKey: make(map[int]string, 1),
				Questions: []Question{
					Question{
						Id: 1,
						QuestionText: "Asset Type",
						DataId: "AssetType",
						Required: true,
						UserResponse: Response{
							IsSelect: true,
							Options: []option{
								option{
									Text: "Roofing",
									Value: "Roofing",
									Selected: true,
								},
								option{
									Text: "Chiller",
									Value: "Chiller",
									Selected: false,
								},
								option{
									Text: "Fan",
									Value: "Fan",
									Selected: false,
								},
							},
						},
					},
				},
				Followups: []Question{
					Question{
						Id: 1,
						QuestionText: "Asset Location Environment",
						DataId: "AssetEnvironment",
						Required: true,
						UserResponse: Response{IsString: true, },
					},
					Question{
						Id: 2,
						QuestionText: "Original Quality or Efficiency",
						DataId: "OriginalEfficiency",
						Required: true,
						UserResponse: Response{IsString: true, },
					},
					Question{
						Id: 3,
						QuestionText: "Annual Run Hour Estimate",
						DataId: "AnnualHours",
						Required: true,
						UserResponse: Response{IsString: true},
					},
				},
			},
			Prompt{
				Id: 4,
				Complete: false,
				FollowupsKey: nil,
				Questions: []Question{
					Question{
						QuestionText: "Asset Size",
						DataId: "AssetSize",
						Required: true,
						UserResponse: Response{
							IsSelect: true,
							Options: []option{
								option{
									Value: "Large",
									Text: "Large",
									Selected: true,
								},
								option{
									Value: "Extra Large",
									Text: "Extra Large",
								},
								option{
									Value: "Small",
									Text: "Small",
								},
								option{
									Value: "Medium",
									Text: "Medium",
								},
							},
						},
					},
				},
			},
			Prompt{
				Id: 5,
				Complete: false,
				FollowupsKey: nil,
				Questions: []Question{
					Question{
						QuestionText: "Maintenance Frequency",
						DataId: "MaintFreq",
						Required: true,
						UserResponse: Response{
							IsSelect: true,
							Options: []option{
								option{
									Value: "Emergency Repair Only",
									Text: "Emergency Repair Only",
								},
								option{
									Value: "Once a Year",
									Text: "Once a Year",
								},
								option{
									Value: "2 to 3 Times a Year",
									Text: "2 to 3 Times a Year",
									Selected: true,
								},
								option{
									Value: "4 or More Times a Year",
									Text: "4 or More Times a Year",
								},
							},
						},
					},
				},
			},
		},
		SelectedPrompt: 0,
	}

	newSurvey.PromptCount = len(newSurvey.Prompts)
	return newSurvey;
}

func (survey Survey)ToBase64() string {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(survey)
	if err != nil {fmt.Println("failed gob encode", err)}
	return base64.StdEncoding.EncodeToString(b.Bytes())
}

func (survey Survey)CheckForFollowups(promptIndex int) ([]Question, int) {
	prompt := survey.Prompts[promptIndex]

	for i, s := range prompt.FollowupsKey {
		if strings.ToUpper(strings.Trim(prompt.Questions[i].UserResponse.Content, " ")) == strings.Trim(s, " ") {
			var follarr = prompt.Followups
			return follarr, len(follarr)
		}
	}

	return nil, 0
}

func FromBase64(str string) Survey {
	m := Survey{}
	by, err := base64.StdEncoding.DecodeString(str)
	if err != nil { fmt.Println("failed base64 gob Decode", err)}
	b:= bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(&m)
	if err != nil { fmt.Println("failed gob decode", err)}
	return m
}




