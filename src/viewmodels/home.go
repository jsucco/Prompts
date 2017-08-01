package viewmodels


import (
	"net/http"
	"models/user"
	_"errors"
)

type Home struct {
	Title string
	User user.Info
}

type Login struct {
	Title string
	Active string
	Member Member
	ErrorMsg string
	HasError bool
}

func GetHome(w http.ResponseWriter, r *http.Request) (Home, error) {

	result := Home{
		Title: "ROINumbers Home",
	}

	//if uo, err := user.GetUserInfo(r); err == nil {
	//	result.User = uo
	//} else {
	//	return Home{}, errors.New("failed to retrieve user info.")
	//}

	return result, nil
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