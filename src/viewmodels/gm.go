package viewmodels

type Gm struct {
	Title string
	Prompts []Control
	SelectedPrompt int
	LastPrompt int
}

type Control struct {
	Label string
	Value string
	DataId string
	Type string
	IsSelect bool
	IsInput bool
	ParentPrompt int
	Options []option
}

type option struct {
	Value string
	Text string
	Selected bool
}

func GetGM() Gm {

	result := Gm{
		Title: "General Questionaire",
		SelectedPrompt: 1,
		LastPrompt: 4,
		Prompts: []Control{
			Control{
				DataId: "StreetAddress",
				Label: "What is your Street Address?",
				Value: "",
				IsInput: true,
				IsSelect: false,
				Type: "text",
				ParentPrompt: 1,
			},
			Control{
				DataId: "ZIP",
				Label: "What is your ZIP code?",
				Value: "",
				IsInput: true,
				IsSelect: false,
				Type: "number",
				ParentPrompt: 1,
			},
			Control{
				DataId: "City",
				Label: "City",
				Value: "",
				IsInput: true,
				IsSelect: false,
				Type: "text",
				ParentPrompt: 1,
			},
			Control{
				DataId: "State",
				Label: "State",
				Value: "",
				IsInput: true,
				IsSelect: false,
				Type: "text",
				ParentPrompt: 1,

			},
			Control{
				DataId: "InstallDate",
				Label: "Asset Install Date",
				Value: "",
				IsInput: true,
				IsSelect: false,
				Type: "date",
				ParentPrompt: 2,
			},
			Control{
				DataId: "AssetSize",
				Label: "Asset Size",
				Value: "Large",
				IsInput: false,
				IsSelect: true,
				Type: "select",
				ParentPrompt: 3,
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
			Control{
				DataId: "MaintFreq",
				Label: "Maintenance Frequency",
				Value: "2 To 3 Times a Year",
				IsInput: false,
				IsSelect: true,
				Type: "select",
				ParentPrompt: 4,
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
	}

	return result
}
