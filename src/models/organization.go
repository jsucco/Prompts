package models

import (
	"cloud.google.com/go/datastore"
	"golang.org/x/net/context"
	"errors"
	"time"
	"crypto/rand"
	"strconv"
)

type Organization struct {
	Name string
	OfficialName string
	Abbreviation string
	EmailAddress string
	Created time.Time
}

func (o *Organization) SaveOrganization() error {
	if len(o.Name) == 0 {
		return errors.New("Organization Name required.")
	}
	if len(o.Abbreviation) == 0 {
		return errors.New("Abbreviated Name required.")
	}

	if len(o.Abbreviation) > 5 {
		return errors.New("Abbreviated Name to long.")
	}

	ctx := context.Background()

	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return err
	}

	kind := "Organization"
	name := o.Abbreviation + time.Now().Month().String() + strconv.Itoa(time.Now().Year()) + "-" + RandStr(12, "alphanum")
	o.Created = time.Now().Local()

	key := datastore.NameKey(kind, name, nil)

	if _, err := client.Put(ctx, key, o); err !=nil {
		return err
	}
	return nil
}

func RandStr(strSize int, randType string) string {

	var dictionary string

	if randType == "alphanum" {
		dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "alpha" {
		dictionary = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "number" {
		dictionary = "0123456789"
	}

	var bytes = make([]byte, strSize)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}