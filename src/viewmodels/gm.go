package viewmodels

type gm struct {
	Title string
	Prompts []prompt
}

type prompt struct {
	label string
	placeholder string
}

func GetGM() gm {

	firstPrompt := prompt{
		label: "Street Address",
		placeholder: "What is your Street Address?",
	}

	secondPrompt := prompt{
		label: "ZIP",
		placeholder: "What is your ZIP code?",
	}

	result := gm{
		Title: "General Questionaire",
		Prompts: []prompt{
			firstPrompt,
			secondPrompt,
		},
	}

	return result
}