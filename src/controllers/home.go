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

func (this *homeController) get(w http.ResponseWriter, r *http.Request) {

	if Authenticated(r) == false {
		http.Redirect(w, r, "/login?returnURL=" + r.RequestURI, 302)
		return
	}

	vm, errh := viewmodels.GetHome(w, r)
	if errh != nil {
		http.Redirect(w, r, "/login?returnURL=" + r.RequestURI, 302)
		return
	}

	err := tpl.ExecuteTemplate(w,"home.gohtml", vm)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (this *homeController) login(w http.ResponseWriter, req *http.Request) {

	w.Header().Add("Content Type", "text/html")
	vm := viewmodels.GetLogin()
	if req.Method == "POST" {

		email := req.FormValue("UserName")
		password := req.FormValue("PassWord")
		if len(password) > 0 {
			//encr_pass, _ := util.Encode(password)

			member, err := models.GetMember(email, password)

			if err == nil {
				SessionId, _ := models.CreateSessionId(member)

				err_s := models.CreateSession(SessionId, member.Id(), member.Email(), member.FirstName(), member.LastName(), member.OrganizationKey(), member.OrganizationName())

				if err_s == nil {

					models.SetSessionCookie(w, SessionId)

					http.Redirect(w, req, "/", 302)
					return

				} else {
					vm.HasError = true;
					vm.ErrorMsg = "get session - " + err_s.Error();
				}
			} else {
				vm.HasError = true;
				vm.ErrorMsg = err.Error();
			}
		} else {
			vm.HasError = true;
			vm.ErrorMsg = "Password length must be greater than 0."
		}

	}

	this.loginTemplate.Execute(w, vm)
}

func Authenticated(req *http.Request) bool {
	_, error := req.Cookie("SessionId")

	if error == nil {
		var cookieval = models.ReadSessionCookie(req)

		if len(cookieval) > 0 {
			_, error := models.GetSession(cookieval)
			if error == nil {
				return true
			}
		}
		return false
	} else {
		return false
	}
}