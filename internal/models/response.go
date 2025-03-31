package models

type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	Country string `json:"country"`
	Lat     string `json:"lat,omitempty"`
	Lng     string `json:"lng,omitempty"`
}

// Маршрут: /api/address/search метод POST
type SearchResponse struct {
	Addresses []*Address `json:"addresses"`
}

// Маршрут: /api/address/geocode метод POST
type GeocodeResponse struct {
	Addresses []*Address `json:"addresses"`
}
