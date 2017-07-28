package models

import (
	"cloud.google.com/go/datastore"
	"golang.org/x/net/context"
	_"errors"
	"time"
	"strconv"
	"strings"
	"errors"
	"models/shard_counter"
)

type Asset struct {
	K *datastore.Key `datastore:"__key__"`
	Name string
	DateAdded time.Time
	DateInstalled time.Time
	Location GeoLocation
	OrganizationKey string
	Index int
}

type GeoLocation struct {
	Latitude float64
	Longitude float64
	StreetNumber int
	Route string
	Locality string
	Administrative_area_level_1 string
	Country string
	PostalCode string
}

func (a *Asset) AddToStore() error {

	ctx := context.Background()

	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return err
	}

	if err = a.ValidateProperties(); err != nil {
		return err
	}

	var ass Asset = Asset{}
	var parent_key = datastore.NameKey("Organization", a.OrganizationKey, nil)

	asset_key := a.Name + "-" + strconv.Itoa(a.Location.StreetNumber) +
		strings.Replace(a.Location.Route, " ", "", -1) +
		a.Location.PostalCode

	var key = datastore.NameKey("Asset", asset_key, parent_key)
	if err := client.Get(ctx, key, &ass); err != datastore.ErrNoSuchEntity {
		return errors.New("Asset already exists at Location. Enter different Name.")
	}

	if ac, errac := shard_counter.Count(ctx, a.OrganizationKey + "-AssetCounter"); errac == nil {
		a.Index = ac
		a.DateAdded = time.Now().Local()
		_, err = client.Put(ctx, key, a)
		if err == nil {
			shard_counter.Increment(ctx, a.OrganizationKey + "-AssetCounter")
		}
	} else {
		err = errac
	}

	return err
}

func (a *Asset)ValidateProperties() error {
	if len(a.Name) == 0 {
		return errors.New("Asset Name required")
	}

	if len(a.Location.Route) == 0 {
		return errors.New("Street Name required")
	}

	if len(a.Location.Locality) == 0 {
		return errors.New("City Name required")
	}

	if a.Location.StreetNumber == 0 {
		return errors.New("Street Number required.")
	}

	if len(a.Location.PostalCode) == 0 {
		return errors.New("Postal Code required.")
	}

	if len(a.OrganizationKey) == 0 {
		return errors.New("Parent Organization key required.")
	}
	return nil
}

func (a *Asset)MapValues(Field string, Value string) {
	switch f := Field; f {
	case "street_number":
		if str_int, err := strconv.Atoi(Value); err == nil {
			a.Location.StreetNumber = str_int
		}
	case "route":
		if len(Value) > 0 {
			a.Location.Route = Value
		}
	case "locality":
		if len(Value) > 0 {
			a.Location.Locality = Value
		}
	case "administrative_area_level_1":
		if len(Value) > 0 {
			a.Location.Administrative_area_level_1 = Value
		}
	case "country":
		if len(Value) > 0 {
			a.Location.Country = Value
		}
	case "postal_code":
		if len(Value) > 0 {
			a.Location.PostalCode = Value
		}
	case "Latitude":
		if flt, err := strconv.ParseFloat(Value, 64); err == nil {
			a.Location.Latitude = flt
		}
	case "Longitude":
		if flt, err := strconv.ParseFloat(Value, 64); err == nil {
			a.Location.Longitude = flt
		}
	case "name":
		if len(Value) > 0 {
			a.Name = Value
		}
	case "InstallDate":
		if len(Value) > 0 {
			if t, err := time.Parse("2006-01-02", Value); err == nil {
				a.DateInstalled = t
			}
		}
	}
}
