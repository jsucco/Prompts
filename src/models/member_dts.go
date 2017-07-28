package models

import (
	"cloud.google.com/go/datastore"
	"golang.org/x/net/context"
	_"time"
	"errors"
	"time"
	"google.golang.org/appengine/memcache"
)

const (
	projectID = "project-alpha-170622"
)

func GetSession(SessionId string) (Session, error) {
	if len(SessionId) == 0 {
		return Session{}, errors.New("SessionId length cannot be 0")
	}

	ctx := context.Background()

	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return Session{}, err
	}
	kind := "Session"
	name := SessionId
	sessionKey := datastore.NameKey(kind, name, nil)
	var session Session
	err = client.Get(ctx, sessionKey, &session)
	if err != nil {
		return Session{}, err
	}
	return session, nil
}

func CreateSession(SessionId string, memberid int, memberemail string, first string, last string, organization_key string, organization_name string) error {
	if len(SessionId) == 0 {
		return errors.New("Blank SessionId not permitted.")
	}

	if len(memberemail) == 0 {
		return errors.New("Blank Email Address not permitted.")
	}

	ctx := context.Background()

	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return err
	}

	kind := "Session"

	name := SessionId

	sessionKey := datastore.NameKey(kind, name, nil)

	new_session := Session{
		SessionId: SessionId,
		MemberId: memberid,
		MemberFirstName: first,
		MemberLastName: last,
		Created: time.Now(),
		OrganizationKey: organization_key,
		OrganizationName: organization_name,
	}

	if _, err := client.Put(ctx, sessionKey, &new_session); err != nil {
		return err
	}

	memcache.JSON.Set(ctx, &memcache.Item{
		Key:        "session-" + new_session.SessionId,
		Object:     &new_session,
		Expiration: 60,
	})

	return nil
}
