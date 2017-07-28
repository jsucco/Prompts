package models

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	_ "time"
	_ "net/mail"
	"time"
	_"bytes"
	"github.com/gorilla/securecookie"
	"net/http"
	"strings"
)

type Member struct {
	id int
	email string
	password string
	first_name string
	last_name string
	parent_organization_id int
	organization_key string
	organization_name string
}

func (this *Member) Id() int {
	return this.id
}

func (this *Member) Email() string {
	return this.email
}

func (this *Member) Password() string {
	return this.password
}

func (this *Member) FirstName() string {
	return this.first_name
}

func (this *Member) LastName() string {
	return this.last_name
}

func (this *Member) ParentOrganizationId() int {
	return this.parent_organization_id
}

func (this * Member) OrganizationKey() string {
	return this.organization_key
}

func (this *Member) OrganizationName() string {
	return this.organization_name
}

func (this *Member) SetId(value int) {
	this.id = value
}

func (this *Member) SetEmail(value string) {
	this.email = value
}

func (this *Member) SetPassword(value string) {
	this.password = value
}

func (this *Member) SetFirstName(value string) {
	this.first_name = value
}

func (this *Member) SetLastName(value string) {
	this.last_name = value
}

func (this *Member) SetParentOrganizationId(value int) {
	this.parent_organization_id = value
}

func (this *Member) SetOrganizationKey(value string) {
	this.organization_key = value
}

func (this *Member) SetOrganizationName(value string) {
	this.organization_name = value
}

type Session struct {
	MemberId int
	MemberFirstName string
	MemberLastName string
	SessionId string
	OrganizationKey string
	OrganizationName string
	Created time.Time
}

func GetMember(email string, password string) (Member, error) {
	db, err := getSessionsConnection()
	defer db.Close()
	if err == nil {
		//pwd := sha256.Sum256([]byte(password))
		//row := db.QueryRow("SELECT id, email, first_name, last_name" +
		//" FROM management.member WHERE UPPER(email) = ? AND password = left(?, 255)", strings.ToUpper(email), hex.EncodeToString(pwd[:]))
		row := db.QueryRow("SELECT member.id, member.email, member.first_name, member.last_name, member.parent_organization_Id, organization.organization_key, organization.name " +
		    "FROM management.member " +
			"JOIN management.organization ON member.parent_organization_Id = management.organization.id " +
			"WHERE UPPER(email) = ? AND password = left(?, 255)", strings.ToUpper(email), password)

		result := Member{}

		err = row.Scan(&result.id, &result.email, &result.first_name, &result.last_name, &result.parent_organization_id, &result.organization_key, &result.organization_name)

		if err == nil {
			return result, nil
		} else {
			return result, err
		}
	} else {
		return Member{}, errors.New("Unable to get database connection")
	}
}


func InsertMember(member Member) (int, error) {
	member_id := 0

	if len(member.email) == 0 {
		return member_id, errors.New("Email Address is required.")
	}

	if len(member.password) == 0 {
		return member_id, errors.New("Password is required.")
	}

	pass_hashed := sha256.Sum256([]byte(member.password))
	//n := len(pass_hashed)
	var pass_hashed_string = hex.EncodeToString(pass_hashed[:])

	if len(member.first_name) == 0 {
		return member_id, errors.New("first name is required")
	}

	if len(member.last_name) == 0 {
		return member_id, errors.New("Last name is required")
	}

	db, err := getSessionsConnection()
	defer db.Close()
	if err == nil {
		//_, err := db.Exec("INSERT INTO session_management.member (email, password, first_name, last_name, parent_organization_id) values('dd', '" + pass_hashed_string + "', 'JS', 'js', 1)")
		_, err := db.Exec("INSERT INTO management.member (email, password, first_name, last_name, parent_organization_id)" +
		"values(?, ?, ?, ?, ?)", member.email, pass_hashed_string, member.first_name, member.last_name, member.parent_organization_id)
		if err == nil {
			return 2, nil
		} else {
			return 0, err
		}
	} else {
		return 0, err
	}
}


func CreateSessionId(member Member) (string, error) {
	result := Session{}
	result.MemberId = member.Id()

	sessionId := sha256.Sum256([]byte(member.Email() + time.Now().Local().String()))

	return hex.EncodeToString(sessionId[:]), nil
}

func GetMemberBySessionId(sessionId string) (Member, error) {
	result := Member{}

	db, err := getSessionsConnection()
	defer db.Close()
	if err == nil {
		err := db.QueryRow("SELECT member.id,email, first_name, last_name, parent_organization_id " +
			"FROM management.session " +
			"JOIN management.member ON member.id = management.session.member_id " +
			"WHERE management.session.session_id = ? " +
			"LIMIT 1", sessionId).Scan(&result.id, &result.email, &result.first_name, &result.last_name, &result.parent_organization_id)

		if err == nil {
			return result, nil
		} else {
			return Member{}, errors.New("Unable to get member for session")
		}
	} else {
		return result, errors.New("Unable to getdatabase connection")
	}
}

func SetSessionCookie(w http.ResponseWriter, value string) {
	var hashkey = []byte(mustGetenv("ROI_HASHKEY"))
	var blockey = []byte(mustGetenv("ROI_BLOCKKEY"))
	var s = securecookie.New(hashkey, blockey)

	content := map[string]string{
		"SessionId":value,
	}
	if encoded, err := s.Encode("SessionId", content); err == nil {
		cookie := &http.Cookie{
			Name: "SessionId",
			Value: encoded,
			Path: "/",
		}
		http.SetCookie(w, cookie)
	}
}

func ReadSessionCookie(r *http.Request) string {
	var hashkey = []byte(mustGetenv("ROI_HASHKEY"))
	var blockey = []byte(mustGetenv("ROI_BLOCKKEY"))
	var s = securecookie.New(hashkey, blockey)

	if cookie, err := r.Cookie("SessionId"); err == nil {
		value := make(map[string]string)
		if err = s.Decode("SessionId", cookie.Value, &value); err == nil {
			return value["SessionId"]
		}
	}
	return ""
}
