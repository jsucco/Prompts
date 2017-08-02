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
	lat float64
	lng float64
}

func (c *Cordinate) SetCordinate(Lat float64, Lng float64) {
	c.lat = Lat
	c.lng = Lng
}

func (r *Result) LoadResults(book []models.Asset) error {
	if len(book) == 0 {
		return errors.New("No assets found.")
	}

	r.Book = book

	for _, a := range book {
		var c = Cordinate{}

		c.SetCordinate(a.Location.Latitude, a.Location.Longitude)
		r.Locations = append(r.Locations, c)
	}

	return nil
}
