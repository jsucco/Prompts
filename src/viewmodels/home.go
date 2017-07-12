package viewmodels


import (

)

type Home struct {
	Title string
	Active string
	Member Member
}

type Login struct {
	Title string
	Active string
	Member Member
	ErrorMsg string
	HasError bool
}

func GetHome() Home {
	result := Home{
		Title: "ROINumbers Home",
		Active: "home",
	}

	return result
}

func GetLogin() Login {
	result := Login{
		Title: "ROINumbers - Login",
		Active: "",
		ErrorMsg: "",
		HasError: false,
	}
	return result
}