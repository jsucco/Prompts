package viewmodels

type gm struct {
	Title string
}

func GetGM() gm {
	result := gm{
		Title: "gm",
	}

	return result
}