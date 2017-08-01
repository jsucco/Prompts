package search

import (
	"models"
	"errors"
)

type Result struct {
	Locations []Cordinate
	Book []models.Asset
}

type Cordinate struct {
	Lat float64
	Lng float64
}

func (r *Result) LoadResults(book []models.Asset) error {
	if len(book) == 0 {
		return errors.New("No assets found.")
	}

	r.Book = book

	for _, a := range book {
		var c Cordinate

		c.Lat = a.Location.Latitude;
		c.Lng = a.Location.Longitude;
		r.Locations = append(r.Locations, c)
	}

	return nil
}
