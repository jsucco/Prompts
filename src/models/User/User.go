package user

import (
	"models"
	"net/http"
	"errors"
)

type Info struct {
	FirstName string
	LastName string
	OrganizationName string
}

func GetUserInfo(req *http.Request) (Info, error) {
	ui := Info{}
	var sessionid = models.ReadSessionCookie(req)
	if len(sessionid) == 0 {
		return Info{}, errors.New("No sessionid was found.")
	}
	sess, err := models.GetUserSession(sessionid, req)
	if err != nil {
		return Info{}, err
	}

	ui.FirstName = sess.MemberFirstName
	ui.LastName = sess.MemberLastName
	ui.OrganizationName = sess.OrganizationName

	return ui, nil
}
