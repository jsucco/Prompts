package models

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	_ "time"
	_ "net/mail"
	"time"
	_"bytes"
	"strings"
	"github.com/gorilla/securecookie"
	"net/http"
)

type Member struct {
	id int
	email string
	password string
	first_name string
	last_name string
	parent_organization_id int
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

type Session struct {
	id int
	memberId int
	sessionId string
}

func (this *Session) Id() int {
	return this.id
}

func (this *Session) MemberId() int {
	return this.memberId
}

func (this *Session) SessionId() string {
	return this.sessionId
}

func (this *Session) SetId(value int) {
	this.id = value
}

func (this *Session) SetMemberId(value int) {
	this.memberId = value
}

func (this *Session) SetSessionId(value string) {
	this.sessionId = value
}

func GetMember(email string, password string) (Member, error) {
	db, err := getSessionsConnection()
	defer db.Close()
	if err == nil {
		//pwd := sha256.Sum256([]byte(password))
		//row := db.QueryRow("SELECT id, email, first_name, last_name" +
		//" FROM management.member WHERE UPPER(email) = ? AND password = left(?, 255)", strings.ToUpper(email), hex.EncodeToString(pwd[:]))
		row := db.QueryRow("SELECT id, email, first_name, last_name" +
			" FROM management.member WHERE UPPER(password) = 'ZJGRROI'")

		result := Member{}

		err = row.Scan(&result.id, &result.email, &result.first_name, &result.last_name)

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

	db, err := getsmDBconnection()

	if err == nil {
		defer db.Close()

		//_, err := db.Exec("INSERT INTO session_management.member (email, password, first_name, last_name, parent_organization_id) values('dd', '" + pass_hashed_string + "', 'JS', 'js', 1)")
		_, err := db.Exec("INSERT INTO session_management.member (email, password, first_name, last_name, parent_organization_id)" +
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

func CreateSession(member Member) (Session, error) {
	result := Session{}
	result.memberId = member.Id()
	sessionId := sha256.Sum256([]byte(member.Email() + time.Now().Format("12:00:00")))
	result.sessionId = hex.EncodeToString(sessionId[:])

	db, err := getsmDBconnection()
	if err == nil {
		defer db.Close()
		res, err := db.Exec("INSERT INTO session_management.session (session_id, member_id)" +
			"VALUES (?, ?);", result.sessionId, member.Id())
		if err == nil {
			id, err := res.LastInsertId()
			if err == nil {
				result.id = int(id)

			}
			return result, nil
		} else {
			return Session{}, err
		}
	} else {
		return Session{}, err
	}
}

func GetMemberBySessionId(sessionId string) (Member, error) {
	result := Member{}

	db, err := getsmDBconnection()
	if err == nil {
		err := db.QueryRow("SELECT member.id, email, first_name, last_name, parent_organization_id " +
			"FROM session_management.session " +
			"JOIN member ON member.id = session.member_id " +
			"WHERE session.session_id = ?", sessionId).Scan(&result.id, &result.email, &result.first_name, &result.last_name, &result.parent_organization_id)

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
	var hashkey = []byte("adksjflk4")
	var blockey = []byte("qwertyuiop123456")
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
	var hashkey = []byte("adksjflk4")
	var blockey = []byte("qwertyuiop123456")
	var s = securecookie.New(hashkey, blockey)

	if cookie, err := r.Cookie("SessionId"); err == nil {
		value := make(map[string]string)
		if err = s.Decode("SessionId", cookie.Value, &value); err == nil {
			return value["SessionId"]
		}
	}
	return ""
}
