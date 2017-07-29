package models

import (
	"fmt"
	"encoding/base64"
	"encoding/gob"
	"bytes"
	"strings"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"errors"
	"time"
	"net/http"
	"strconv"
)

type Survey struct {
	Title string
	OrganizationId int
	Finished time.Time
	Updated time.Time
	Type string
	Prompts []Prompt
	PromptCount int
	Completed bool
	SelectedPrompt int
	LastPrompt int
	K *datastore.Key `datastore:"__key__"`
}

type AssetSurvey struct {
	Title string
	OrganizationId int
	Finished time.Time
	Updated time.Time
	Type string
	Content []byte
	PromptCount int
	Completed bool
	SelectedPrompt int
	LastPrompt int
	K *datastore.Key `datastore:"__key__"`
}

var (
	UserRequest *http.Request
	NewAsset Asset
)

func (s *Survey) MapAllResponses(req *http.Request) error {
	if req.Method == "POST" {
		s.Title = "Asset Assessment"

		if len(s.Prompts) > 0 {
			NewAsset = Asset{}

			for i, _ := range s.Prompts {
				for j, _ := range s.Prompts[i].Questions {

					s.Prompts[i].Questions[j], _ = s.Prompts[i].Questions[j].MapQuestion(req)
				}
			}
		} else {
			errors.New("survey did not contain any prompts.")
		}
	} else {
		return errors.New("method only maps inputs from a post request.")
	}
	return nil
}

func (q *Question) MapQuestion(req *http.Request) (Question, error) {
	var inputval = req.FormValue(q.DataId)

	if len(inputval) > 0 {
		q.UserResponse.Content = inputval
		NewAsset.MapValues(q.DataId, inputval)
	}
	if len(q.Followups) > 0 {
		for i, _ := range q.Followups {
			fq, err := q.Followups[i].MapFwQuestion(req)
			if err == nil {
				q.Followups[i] = fq
			}
		}
	}
	return *q, nil
}

func (fq *FwQuestion) MapFwQuestion(req *http.Request) (FwQuestion, error) {
	var inputval = req.FormValue(fq.DataId)

	if len(inputval) > 0 {
		fq.UserResponse.Content = inputval
		NewAsset.MapValues(fq.DataId, inputval)
	}
	return *fq, nil
}

func (s *Survey) LoadSurvey(OrganizationKeyStr string, SurveyKeyStr string, req *http.Request) error {
	ctx := appengine.NewContext(req)

	parent_kind := "Organization"
	kind := "AssetSurvey"

	key := datastore.NewKey(ctx, parent_kind, OrganizationKeyStr, 0, nil)

	survey_key := datastore.NewKey(ctx, kind, SurveyKeyStr, 0, key)
	as := AssetSurvey{}

	err := datastore.Get(ctx, survey_key, &as)
	if err != nil {
		return err
	}
	if len(as.Content) == 0 {
		return errors.New("No content found in survey store object.")
	}
	ns, errn := DeserializeBuffer(string(as.Content))
	if errn != false {
		return errors.New("Load failed.")
	}
	s = &ns
	return nil
}

func (s *Survey) SaveSurvey(OrganizationKey string, req *http.Request) error {
	if len(s.Prompts) == 0 {
		return errors.New("Survey must contain at least one prompts.")
	}

	if len(s.Prompts[0].Questions) == 0 {
		return errors.New("Survey must contain at least one question.")
	}

	if len(OrganizationKey) == 0 {
		return errors.New("OrganizationKey required in order to Save Survey.")
	}
	as := AssetSurvey{}
	as.Content = []byte(s.ToBase64())
	as.Title = s.Title
	as.OrganizationId = s.OrganizationId
	as.Finished = s.Finished
	as.Updated = time.Now().Local()
	as.Type = s.Type
	as.PromptCount = s.PromptCount
	as.Completed = s.Completed
	as.SelectedPrompt = s.SelectedPrompt
	as.LastPrompt = s.LastPrompt
	as.Completed = true

	ctx := appengine.NewContext(req)

	kind := "AssetSurvey"
	name := s.Type + time.Now().Month().String() + strconv.Itoa(time.Now().Year()) + "-" + RandStr(12, "alphanum")

	parent_key := datastore.NewKey(ctx,"Organization", OrganizationKey,0, nil)

	survey_key := datastore.NewKey(ctx, kind, name, 0,parent_key)

	if _, err := datastore.Put(ctx, survey_key, &as); err != nil {
		return err
	}
	return nil
}

func (s *Survey) AddNewAsset(OrganizationKey string, req *http.Request) error {
	if len(OrganizationKey) == 0 {
		return errors.New("OrganizationKey required.")
	}

	NewAsset.OrganizationKey = OrganizationKey

	if err := NewAsset.AddToStore(req); err != nil {
		return err
	}
	return nil
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
	Followups []FwQuestion
	FollowupAddress []string
	Permissions []string
}

type FwQuestion struct {
	Id int
	Address string
	SurveyLevel int
	QuestionText string
	DataId string
	UserResponse Response
	Required bool
	ErrorMessage string
	ErrorFlag bool
	Permissions []string
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
			var fu = list_followups(q.Followups)
			qr = append(qr, fu...)
		}

	}
	return qr
}

func list_followups(qs []FwQuestion) []Question {
	qr := []Question{}

	for _, q := range qs {
		qn := Question{
			Id: q.Id,
			Address: q.Address,
			SurveyLevel: q.SurveyLevel,
			QuestionText: q.QuestionText,
			DataId: q.DataId,
			Required: q.Required,
			UserResponse: q.UserResponse,
			Permissions: q.Permissions,
		}
		qr = append(qr, qn)
	}
	return qr
}

func AssembleGM() Survey {

	newSurvey := Survey{
		Title: "Asset Assessment",
		Type: "Asset",
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
						Followups: []FwQuestion{
							FwQuestion{
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
					Question{
						Id: 8,
						Address: "8",
						SurveyLevel: 1,
						QuestionText: "Latitude",
						DataId: "Latitude",
						Required: true,
						UserResponse: Response{
							Default: Input{
								Type: "number",
							},
						},
					},
					Question{
						Id: 9,
						Address: "9",
						SurveyLevel: 1,
						QuestionText: "Longitude",
						DataId: "Longitude",
						Required: true,
						UserResponse: Response{
							Default: Input{
								Type: "number",
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
						Address: "10",
						SurveyLevel: 1,
						QuestionText: "Asset Name",
						DataId: "name",
						Required: true,
						UserResponse: Response{
							Default: Input{
								Type: "text",
							},
						},
					},
					Question{
						Id: 2,
						Address: "11",
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
						Id: 3,
						Address: "12",
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
							"12_1",
							"12_2",
							"12_3",
						},
						Followups: []FwQuestion {
							FwQuestion{
								Id: 1,
								Address: "12_1",
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
							FwQuestion{
								Id: 2,
								Address: "12_2",
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
							FwQuestion{
								Id: 3,
								Address: "12_3",
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
						Id: 4,
						Address: "13",
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
								RefDataId: "Asset_Use",
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
						Address: "14",
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
							"14_1",
							"14_2",
							"14_3",
						},
						Followups: []FwQuestion{
							FwQuestion{
								Id: 1,
								Address: "14_1",
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
							FwQuestion{
								Id: 2,
								Address: "14_2",
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
							FwQuestion{
								Id: 3,
								Address: "14_3",
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
						Address: "15",
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
						Address: "16",
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




