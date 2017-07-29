package models

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"time"
	"net/http"
	"errors"
)

func GetUserSession(SessionId string, req *http.Request) (Session, error) {
	if len(SessionId) == 0 {
		return Session{}, errors.New("SessionId length cannot be 0")
	}

	ctx := appengine.NewContext(req)

	kind := "Session"
	name := SessionId
	sessionKey := datastore.NewKey(ctx, kind, name, 0, nil)
	var session Session
	err := datastore.Get(ctx, sessionKey, &session)
	if err != nil {
		return Session{}, err
	}
	return session, nil
}

func CreateUserSession(SessionId string, memberid int, memberemail string, first string, last string, organization_key string, organization_name string, req *http.Request) error {
	if len(SessionId) == 0 {
		return errors.New("Blank SessionId not permitted.")
	}

	if len(memberemail) == 0 {
		return errors.New("Blank Email Address not permitted.")
	}

	ctx := appengine.NewContext(req)

	kind := "Session"

	name := SessionId

	sessionKey := datastore.NewKey(ctx, kind, name, 0, nil)

	new_session := Session{
		SessionId: SessionId,
		MemberId: memberid,
		MemberFirstName: first,
		MemberLastName: last,
		Created: time.Now(),
		OrganizationKey: organization_key,
		OrganizationName: organization_name,
	}

	if _, err := datastore.Put(ctx, sessionKey, &new_session); err != nil {
		return err
	}

	//memcache.JSON.Set(ctx, &memcache.Item{
	//	Key:        "session-" + new_session.SessionId,
	//	Object:     &new_session,
	//	Expiration: 60,
	//})

	return nil
}
