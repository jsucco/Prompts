package controllers

import (
	"net/http"
	"viewmodels"
	"text/template"
	"models"
)

type homeController struct {
	template *template.Template
	loginTemplate *template.Template
}

func (this *homeController) get(w http.ResponseWriter, req *http.Request) {

	if Authenticated(req) == false {
		http.Redirect(w, req, "/login?returnURL=" + req.RequestURI, 302)
		return
	}

	vm := viewmodels.GetHome()

	w.Header().Add("Content Type", "text/html")

	if Authenticated(req) == true {
		this.template.Execute(w, vm)
	} else {
		http.Redirect(w, req, "/login?returnURL=" + req.RequestURI, 302)
	}

}

func (this *homeController) login(w http.ResponseWriter, req *http.Request) {

	w.Header().Add("Content Type", "text/html")
	vm := viewmodels.GetLogin()
	if req.Method == "POST" {

		email := req.FormValue("UserName")
		password := req.FormValue("PassWord")

		member, err := models.GetMember(email, password)

		if err == nil {
			session, err_s := models.CreateSession(member)
			if err_s == nil {

				models.SetSessionCookie(w, session.SessionId())

				http.Redirect(w, req, "/survey", 302)
				return

			} else {
				vm.HasError = true;
				vm.ErrorMsg = "get session - " + err_s.Error();
			}
		} else {
			vm.HasError = true;
			vm.ErrorMsg = err.Error();
		}
	}

	this.loginTemplate.Execute(w, vm)
}

func Authenticated(req *http.Request) bool {
	_, error := req.Cookie("SessionId")

	if error == nil {
		var cookieval = models.ReadSessionCookie(req)

		if len(cookieval) > 0 {
			_, error := models.GetMemberBySessionId(cookieval)
			if error == nil {
				return true
			}
		}
		return false
	} else {
		return false
	}
}