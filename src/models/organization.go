package models

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"errors"
	"time"
	"crypto/rand"
	"strconv"
	"net/http"
)

type Organization struct {
	Name string
	OfficialName string
	Abbreviation string
	EmailAddress string
	Created time.Time
	Key string
}

func (o *Organization) SaveOrganization(req *http.Request) error {
	if len(o.Name) == 0 {
		return errors.New("Organization Name required.")
	}
	if len(o.Abbreviation) == 0 {
		return errors.New("Abbreviated Name required.")
	}

	if len(o.Abbreviation) > 5 {
		return errors.New("Abbreviated Name to long.")
	}

	ctx := appengine.NewContext(req)

	kind := "Organization"
	name := o.Abbreviation + time.Now().Month().String() + strconv.Itoa(time.Now().Year()) + "-" + RandStr(12, "alphanum")
	o.Created = time.Now().Local()

	key := datastore.NewKey(ctx, kind, name, 0, nil)

	if _, err := datastore.Put(ctx, key, o); err !=nil {
		return err
	}
	return nil
}

func (o *Organization) GetAssets(req *http.Request) ([]Asset, error) {
	if len(o.Key) == 0 {
		return []Asset{}, errors.New("Organization Key required.")
	}

	ctx := appengine.NewContext(req)

	parentkey := datastore.NewKey(ctx, "Organization", o.Key, 0, nil)

	query := datastore.NewQuery("Asset").Ancestor(parentkey)

	var book []Asset

	for t := query.Run(ctx); ; {
		var a Asset
		_, err := t.Next(&a)
		if err == datastore.Done {
			break
		}
		if err != nil {
			return book, err
		}
		book = append(book, a)
	}
	return book, nil
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