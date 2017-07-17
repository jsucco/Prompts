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
	Complete bool
	SelectedPrompt int
	Title string
}

type Question struct {
	Id int
	Address string
	SurveyLevel int
	QuestionText string
	DataId string
	UserResponse Response
	Required bool
	ErrorMessage string
	ErrorFlag bool
	Followups []Question
	Permissions []string
	FollowupAddress []string
}

type Response struct {
	Content string
	Options []option
	Complete bool
	IsSelect bool
	IsMultiLine bool
	IsSlider bool
	Slide Slider
	Default Input
}

type Slider struct {
	Min int
	Max int
	Interval int
	LabelText string
	AltLabelText string
	RefDataId string
}

type Input struct {
	Type string
	AutoComplete bool
	AutoCompType string
}

type option struct {
	Value string
	Text string
	RefA int
	Selected bool
}

func DeserializeBuffer(buffer string) (Survey, bool) {
	if len(strings.Trim(buffer, " ")) == 0 {
		return Survey{}, false
	}
	return FromBase64(buffer), true
}

func (p *Prompt)ListAllQuestions() []Question {
	if len(p.Questions) > 0 {
		return list_questions(p.Questions)
	}
	return []Question{}
}

func list_questions(qs []Question) []Question {
	qr := []Question{}

	for _, q := range qs {
		qr = append(qr, q)

		if len(q.Followups) > 0 {
			var fu = list_questions(q.Followups)
			qr = append(qr, fu...)
		}

	}
	return qr
}

func AssembleGM() Survey {

	newSurvey := Survey{
		Title: "General Manager",
		Prompts: []Prompt{
			Prompt{
				Id: 1,
				Complete: false,
				Title: "Location",
				Questions: []Question{
					Question{
						Id: 1,
						Address: "1",
						SurveyLevel: 1,
						QuestionText: "What is your Street Address?",
						DataId: "autocomplete",
						Required: true,
						UserResponse: Response{
							Default: Input{
								Type:"text",
								AutoComplete: true,
							},
						},
					},
					Question{
						Id: 2,
						Address: "2",
						SurveyLevel: 1,
						QuestionText: "Street Number",
						DataId: "street_number",
						Required: true,
						UserResponse: Response{
							Default: Input{
								Type: "text",
							},
						},
					},
					Question{
						Id: 3,
						Address: "3",
						SurveyLevel: 1,
						QuestionText: "Street Name",
						DataId: "route",
						Required: true,
						UserResponse: Response{
							Default: Input{
								Type: "text",
							},
						},
					},
					Question{
						Id: 4,
						Address: "4",
						SurveyLevel: 1,
						QuestionText: "City",
						DataId: "locality",
						Required: true,
						UserResponse: Response{
							Default: Input{
								Type: "text",
							},
						},
						FollowupAddress: []string{
							"3_1",
						},
						Followups: []Question{
							Question{
								Id: 1,
								Address: "3_1",
								SurveyLevel: 2,
								QuestionText: "Please describe the condition of the City?",
								DataId: "city_condition",
								Required: true,
								UserResponse: Response{IsMultiLine: true },
								Permissions:[]string{
									"Seattle",
									"Cincinnati",
								},
							},
						},
					},
					Question{
						Id: 5,
						Address: "5",
						SurveyLevel: 1,
						QuestionText: "State",
						DataId: "administrative_area_level_1",
						Required: true,
						UserResponse: Response{
							Default: Input{
								Type: "text",
							},
						},
					},
					Question{
						Id: 6,
						Address: "6",
						SurveyLevel: 1,
						QuestionText: "Country",
						DataId: "country",
						Required: true,
						UserResponse: Response{
							Default: Input{
								Type: "text",
							},
						},
					},
					Question{
						Id: 7,
						Address: "7",
						SurveyLevel: 1,
						QuestionText: "ZIP",
						DataId: "postal_code",
						Required: true,
						UserResponse: Response{
							Default: Input{
								Type: "text",
							},
						},
					},
				},
			},
			Prompt{
				Id: 2,
				Complete: false,
				Title: "USE GROUP",
				Questions: []Question{
					Question{
						Id: 1,
						Address: "5",
						SurveyLevel: 1,
						QuestionText: "Asset Install Date",
						DataId: "InstallDate",
						Required: true,
						UserResponse: Response{
							Default: Input{
								Type: "date",
							},
						},
					},
					Question{
						Id: 2,
						Address: "6",
						SurveyLevel: 1,
						QuestionText: "Asset Use",
						DataId: "Asset_Use",
						Required: true,
						UserResponse: Response{
							IsSelect: true,
							Options: []option{
								option{
									Text: "Office",
									Value: "400",
								},
								option{
									Text: "Hotel",
									Value: "350",
								},
								option{
									Text: "Restaurant",
									Value: "300",
								},
								option{
									Text: "Industrial",
									Value: "400",
								},
								option{
									Text: "Warehouse",
									Value: "400",
								},
								option{
									Text: "Residence",
									Value: "500",
								},
								option{
									Text: "Educational",
									Value: "400",
								},
								option{
									Text: "Retail",
									Value: "450",
								},
								option{
									Text: "Worship",
									Value: "350",
								},
								option{
									Text: "Garage",
									Value: "Garage",
								},
								option{
									Text: "Hospital",
									Value: "250",
								},
								option{
									Text: "Casino",
									Value: "350",
								},
							},
						},
						FollowupAddress: []string{
							"6_1",
							"6_2",
							"6_3",
						},
						Followups: []Question {
							Question{
								Id: 1,
								Address: "6_1",
								SurveyLevel: 2,
								QuestionText: "Saturday Occupancy",
								DataId: "Saturday_Occupancy",
								Required: true,
								UserResponse: Response{
									Default: Input{
										Type: "checkBox",
									},
								},
								Permissions: []string{
									"Office",
								},
							},
							Question{
								Id: 2,
								Address: "6_2",
								SurveyLevel: 2,
								QuestionText: "Sunday Occupancy",
								DataId: "Sunday_Occupancy",
								Required: true,
								UserResponse: Response{
									Default: Input{
										Type: "checkBox",
									},
								},
								Permissions: []string{
									"Office",
								},
							},
							Question{
								Id: 3,
								Address: "6_3",
								SurveyLevel: 2,
								QuestionText: "Major Holiday Occupancy",
								DataId: "Major_Holiday_Occupancy",
								Required: true,
								UserResponse: Response{
									Default: Input{
										Type: "checkBox",
									},
								},
								Permissions: []string{
									"Office",
								},
							},
						},
					},
					Question{
						Id: 3,
						Address: "7",
						SurveyLevel: 1,
						QuestionText: "Asset Size",
						DataId: "Size_Asset",
						Required: true,
						UserResponse: Response{
							IsSlider: true,
							Slide: Slider{
								Min: 0,
								Max: 10000,
								Interval: 1000,
								LabelText: "Building Square Ft.",
								AltLabelText: "Chiller Ton(s)",
								RefDataId: "Use_Group",
							},
						},
					},
				},
			},
			Prompt{
				Id: 3,
				Complete: false,
				Questions: []Question{
					Question{
						Id: 1,
						Address: "6",
						SurveyLevel: 1,
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
						FollowupAddress: []string{
							"6_1",
							"6_2",
							"6_3",
						},
						Followups: []Question{
							Question{
								Id: 1,
								Address: "6_1",
								SurveyLevel: 2,
								QuestionText: "Asset Location Environment",
								DataId: "AssetEnvironment",
								Required: true,
								UserResponse: Response{
									Default: Input{
										Type: "text",
									},
								},
								Permissions: []string{
									"Chiller",
								},
							},
							Question{
								Id: 2,
								Address: "6_2",
								SurveyLevel: 2,
								QuestionText: "Original Quality or Efficiency",
								DataId: "OriginalEfficiency",
								Required: true,
								UserResponse: Response{
									Default: Input{
										Type: "number",
									},
								},
								Permissions: []string{
									"Chiller",
								},
							},
							Question{
								Id: 3,
								Address: "6_3",
								SurveyLevel: 2,
								QuestionText: "Annual Run Hour Estimate",
								DataId: "AnnualHours",
								Required: true,
								UserResponse: Response{
									Default: Input{
										Type: "number",
									},
								},
								Permissions: []string{
									"Chiller",
								},
							},
						},
					},
				},
			},
			Prompt{
				Id: 4,
				Complete: false,
				Questions: []Question{
					Question{
						Id: 1,
						Address: "7",
						SurveyLevel: 1,
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
				Questions: []Question{
					Question{
						Id: 1,
						Address: "8",
						SurveyLevel: 1,
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

//func (survey Survey)CheckForFollowups(promptIndex int) ([]Question, int) {
//	prompt := survey.Prompts[promptIndex]
//
//	for i, s := range prompt.FollowupsKey {
//		if strings.ToUpper(strings.Trim(prompt.Questions[i].UserResponse.Content, " ")) == strings.Trim(s, " ") {
//			var follarr = prompt.Followups
//			return follarr, len(follarr)
//		}
//	}
//
//	return nil, 0
//}

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




