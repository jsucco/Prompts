package viewmodels


import (

)

type Home struct {
	Title string
	Active string
}

func GetHome() Home {
	result := Home{
		Title: "ROINumbers Home",
		Active: "home",
	}

	return result
}