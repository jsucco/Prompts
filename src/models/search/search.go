package search

import (
	"models"
	"errors"
	"strconv"
)

type Result struct {
	Locations []Cordinate
	Book []models.Asset
}

type Cordinate struct {
	lat string
	lng string
}

func (r *Result) LoadResults(book []models.Asset) error {
	if len(book) == 0 {
		return errors.New("No assets found.")
	}

	r.Book = book

	for _, a := range book {
		c := Cordinate{
			lat: "t" + strconv.FormatFloat(a.Location.Latitude, 'f', 6, 64),
			lng: "t" + strconv.FormatFloat(a.Location.Longitude, 'f', 6, 64),
		}
		r.Locations = append(r.Locations, c)
	}

	return nil
}
